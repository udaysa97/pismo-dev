package otpservice

import (
	"context"
	"encoding/base32"
	"fmt"
	"math/rand"
	apiLogger "pismo-dev/api/logger"
	encryption "pismo-dev/commonpkg/encryptionutils"
	"pismo-dev/constants"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/cache"
	"pismo-dev/pkg/logger"
	"strconv"
	"time"

	"github.com/pquerna/otp/totp"
)

type OTPSvc struct {
	ServiceName string
	countCacheW *cache.CacheWrapper[string, int]
	hashCacheW  *cache.CacheWrapper[string, string]
	services    RequiredServices
}

func NewOTPSvc(cacheWstring *cache.CacheWrapper[string, string], cacheWint *cache.CacheWrapper[string, int]) *OTPSvc {
	return &OTPSvc{
		ServiceName: "OTPService",
		countCacheW: cacheWint,
		hashCacheW:  cacheWstring,
	}
}

func (svc *OTPSvc) SetRequiredServices(services RequiredServices) {
	svc.services = services
}

func (svc *OTPSvc) GenerateOTP(ctx context.Context, user types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface) (bool, error) {
	logger.Info("NFT Transfer Request", map[string]interface{}{"context": ctx, "user": user.Id, "NFTDetails": nftTokenDetails})

	if userEligible := svc.CheckEligibility(ctx, user); !userEligible {
		return false, fmt.Errorf("Unauthorized")
	}

	if matched, err := svc.MatchReloginPin(ctx, user); !matched {
		return false, err
	}
	// In future can add logic to check token/NFT balance of user here

	// //SEND OTP
	var err error
	user, err = svc.requestTransferOTP(ctx, user, constants.NFT_TRANSFER_PURPOSE)
	if err != nil {
		return false, err
	}
	if sent := svc.sendMail(ctx, user, constants.NFT_TRANSFER_PURPOSE); !sent && appconfig.ENV != constants.DEVELOPMENT_ENV {
		logger.Error("SendOTPMail:could not Generate OTP", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "user": user.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNPROCESSABLE_ENTITY_ERROR].ErrorCode, "Error sending OTP")
		return false, fmt.Errorf("could not Generate OTP")
	}
	return true, nil
}

func (svc *OTPSvc) sendMail(ctx context.Context, user types.UserDetailsInterface, purpose string) bool {
	via := "smtp"
	if !appconfig.IS_PRODUCTION {
		via = "ses"
	}

	if len(user.UserOTP) == 0 {
		logger.Error("SendOTPMail: Could not retreive OTP of user from object", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": user.Id})
		return false
	}

	emailApiObject := map[string]interface{}{
		"coindcx_id": user.Id,
		//Todo, to be revamped after decoupling of communication service
		"template_name": fmt.Sprintf("%s%s", constants.COMMUNICATION_SERVICE_PURPOSE_PREFIX, constants.NFT_TRANSFER_PURPOSE), //TODO: constants.COMMUNICATION_SERVICE_PURPOSE_PREFIX + purpose,
		"source":        constants.COMMUNICATION_SERVICE_SOURCE,
		"via":           via,
		"asynchronous":  true,
	}

	payload := map[string]interface{}{
		"extras": map[string]interface{}{
			"emailPurpose": "NFT Transfer",
		},
		"token": user.UserOTP,
	}

	deviceDetails := map[string]interface{}{
		"source":     user.DeviceDetails.Source,
		"ip_address": user.DeviceDetails.IPAddress,
		"device":     user.DeviceDetails.Device,
	}
	emailApiObject["payload"] = payload
	emailApiObject["device_details"] = deviceDetails
	emailPayload := map[string]interface{}{
		"template": emailApiObject,
	}
	if sent, err := svc.services.EmailSvc.SendMail(ctx, emailPayload); !sent {
		logger.Error("EmailSvc: Error sending mail", map[string]interface{}{"context": ctx, "error": err, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": user.Id})
		return false
	}
	return true
}

func (svc *OTPSvc) requestTransferOTP(ctx context.Context, user types.UserDetailsInterface, purpose string) (types.UserDetailsInterface, error) {
	var err error
	user.UserOTP, err = svc.generateOTP(ctx, user.Id)
	if err != nil {
		return user, err
	}
	encrypted, err := encryption.HashString(user.UserOTP)
	if err != nil {
		logger.Error("RequestTransferOTP: Could not encrypt OTP", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err.Error(), "user": user.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Could not generate hash for otp")
		return user, err
	}
	err = svc.hashCacheW.Driver.SetEx(user.Id, encrypted, purpose, 5*time.Minute)
	if err != nil {
		logger.Error("RequestTransferOTP: Could not Save OTP in cache", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "error": err.Error(), "user": user.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Could not save OTP")
		return user, err
	}
	return user, nil
}

func (svc *OTPSvc) generateOTP(ctx context.Context, secretSeed string) (string, error) {
	magicOTP := constants.MAGIC_OTP["testing"]
	if len(magicOTP) > 0 && !appconfig.IS_PRODUCTION {
		return magicOTP, nil
	}

	secret := secretSeed + strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
	secretBase32 := base32.StdEncoding.EncodeToString([]byte(secret))
	otp, err := totp.GenerateCode(secretBase32, time.Now())
	if err != nil {
		logger.Error("OTP: Could not generate OTP", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "userId": secretSeed, "secret": secret, "error": err.Error()})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Could not generate otp")
		return "", err
	}
	return otp, nil
}

func (svc *OTPSvc) CheckEligibility(ctx context.Context, user types.UserDetailsInterface) bool {

	var coolOffCount, logOutCount int
	var err error
	var found bool

	coolOffCountTTL := appconfig.DEFI_COOL_OFF_PERIOD

	coolOffCountKey := svc.coolOffCountKey(user.Id)
	logOutCountKey := svc.logOutCountKey(user.Id)

	if coolOffCount, found, err = svc.countCacheW.Driver.Get(coolOffCountKey, ""); err != nil {
		logger.Error("CheckUserEligibility, Issue while retreiving cache key for cooloffcount", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "found": found, "error": err, "Result": coolOffCount})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Could not retrieve cool off count")
		return false
	} else if !found {
		coolOffCount = 0
	}

	if logOutCount, found, err = svc.countCacheW.Driver.Get(logOutCountKey, ""); err != nil {
		logger.Error("CheckUserEligibility, Issue while retreiving cache key for logoffcount", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "found": found, "error": err, "Result": logOutCount})
		apiLogger.GinLogErrorAndAbort(ctx, 401, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Could not retrieve cool off count")
		return false
	} else if !found {
		logOutCount = 0
	}

	logger.Debug(fmt.Sprintf("COOL OFF COUNT : %d", coolOffCount))
	logger.Debug(fmt.Sprintf("LOG OFF COUNT : %d", logOutCount))
	// //FORCE LOGOUT CHECK
	if logOutCount >= constants.LOGOUT_COUNT {
		svc.countCacheW.Driver.Delete(logOutCountKey, "")
		//Make auth service call to force logout
		svc.services.AuthSvc.ForceLogout(ctx, user)
		logger.Error("Multiple wrong PIN attempts in logoutcount", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.UNAUTHORISED_ERROR].ErrorCode, "user": user.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.UNAUTHORISED_ERROR].HttpStatus, constants.ERROR_TYPES[constants.UNAUTHORISED_ERROR].ErrorCode, "Multiple wrong PIN attempts in logoutcount")
		return false
	}
	if coolOffCount >= constants.COOL_OFF_COUNT {
		svc.countCacheW.Driver.Expire(coolOffCountKey, "", time.Duration(coolOffCountTTL)*time.Minute)
		logger.Error("Multiple wrong PIN attempts in cooloff count", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].ErrorCode, "user": user.Id})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].HttpStatus, constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].ErrorCode, "Multiple wrong PIN attempts in cooloff count")
		return false
	}
	return true

}

func (svc *OTPSvc) MatchReloginPin(ctx context.Context, user types.UserDetailsInterface) (bool, error) {
	if len(user.ReloginPin) == 0 {
		logger.Error("MatchReloginPin: No pin set in request", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].ErrorCode, "user": user})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].HttpStatus, constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].ErrorCode, "Invalid Relogin pin")
		return false, fmt.Errorf("no Login PIN passed")
	}

	if pinVerified, err := svc.services.AuthSvc.VerifyReloginPin(ctx, user); !pinVerified || err != nil {
		if !svc.onInvalidPin(ctx, user) {
			return false, fmt.Errorf("Issue while checking pin cache and eligibility")
		}
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].HttpStatus, constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].ErrorCode, "Invalid login Pin")
		return false, fmt.Errorf("invalid login Pin")
	}

	svc.countCacheW.Driver.Delete(svc.coolOffCountKey(user.Id), "")
	svc.countCacheW.Driver.Delete(svc.logOutCountKey(user.Id), "")
	return true, nil
}

func (svc *OTPSvc) onInvalidPin(ctx context.Context, user types.UserDetailsInterface) bool {
	var logOutCount int
	var found bool
	var err error
	coolOffCountKey := svc.coolOffCountKey(user.Id)
	svc.countCacheW.Driver.Incr(coolOffCountKey, "")

	logOutCountKey := svc.logOutCountKey(user.Id)
	if logOutCount, found, err = svc.countCacheW.Driver.Get(logOutCountKey, ""); err != nil {
		logger.Error("CheckUserEligibility, Issue while retreiving cache key for logoffcount", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].ErrorCode, "found": found, "error": err, "Result": logOutCount})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].HttpStatus, constants.ERROR_TYPES[constants.INVALID_PIN_ERROR].ErrorCode, "Could not retrieve logout count")
		return false
	} else if !found {
		logOutCount = 0
	}

	logOutCountTTL := appconfig.DEFI_FORCE_LOGOUT_PERIOD
	svc.countCacheW.Driver.Incr(logOutCountKey, "")
	if logOutCount == 1 {
		svc.countCacheW.Driver.Expire(logOutCountKey, "", time.Duration(logOutCountTTL)*time.Minute)
	}

	return svc.CheckEligibility(ctx, user)
}

func (svc *OTPSvc) MatchTransferPin(ctx context.Context, user types.UserDetailsInterface) bool {
	acquired, err := svc.hashCacheW.Driver.Mutex(user.UserOTP+"_"+user.Id, constants.NFT_TRANSFER_PURPOSE, "1", 5*time.Second)
	if err != nil || !acquired {
		logger.Error("MatchTransferPin: Failed to acquire lock", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": user.Id, "error": err.Error(), "acquired": acquired})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.PROCESS_ERROR].HttpStatus, constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "Please enter Correct OTP or OTP is Invalid")
		return false
	}
	userOtp, found, err := svc.hashCacheW.Driver.Get(user.Id, constants.NFT_TRANSFER_PURPOSE)
	if err != nil {
		logger.Error("MatchTransferPin: Could not retreive saved OTP", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].ErrorCode, "user": user.Id, "error": err.Error(), "Otpfound": found})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].HttpStatus, constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].ErrorCode, "Please enter Correct OTP or OTP is Invalid")
		return false
	}
	if !found {
		logger.Error("MatchTransferPin: User OTP not found", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].ErrorCode, "user": user.Id, "Otpfound": found})
		apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].HttpStatus, constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].ErrorCode, "Please enter Correct OTP or OTP is Invalid")
		return false

	}
	err = svc.hashCacheW.Driver.Delete(user.UserOTP+"_"+user.Id, constants.NFT_TRANSFER_PURPOSE)
	if err != nil {
		logger.Error("MatchTransferPin: Error deleting mutex key", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.PROCESS_ERROR].ErrorCode, "user": user.Id, "error": err.Error()})
	}
	if encryption.ValidateString(userOtp, user.UserOTP) {
		err = svc.hashCacheW.Driver.Delete(user.Id, constants.NFT_TRANSFER_PURPOSE)
		if err != nil {
			logger.Error("Could not delete User OTP", map[string]interface{}{"context": ctx, "error": err.Error(), "userId": user.Id})
		}
		return true
	}
	logger.Error("MatchTransferPin: OTP Does not match", map[string]interface{}{"context": ctx, "errCode": constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].ErrorCode, "user": user.Id})
	apiLogger.GinLogErrorAndAbort(ctx, constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].HttpStatus, constants.ERROR_TYPES[constants.INVALID_OTP_ERROR].ErrorCode, "Please enter Correct OTP or OTP is Invalid")
	return false

}

func (svc *OTPSvc) coolOffCountKey(userId string) string {
	return fmt.Sprintf("%s::cool_off_count", userId)
}

func (svc *OTPSvc) logOutCountKey(userId string) string {
	return fmt.Sprintf("%s::log_off_count", userId)
}
