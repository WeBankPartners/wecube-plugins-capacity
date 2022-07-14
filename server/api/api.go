package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/api/v1/capacity"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/util/log"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func InitHttpServer(port int) {
	http.Handle("/capacity/api/v1/monitor/search", handleWithLog(capacity.MonitorSearchHandler))
	http.Handle("/capacity/api/v1/monitor/chart", handleWithLog(capacity.MonitorDataHandler))
	http.Handle("/capacity/api/v1/r/data", handleWithLog(capacity.RJustifyDataHandler))
	http.Handle("/capacity/api/v1/r/analyze", handleWithLog(capacity.RAnalyzeHandler))
	http.Handle("/capacity/api/v1/r/calc", handleWithLog(capacity.RCalcDataHandle))
	http.Handle("/capacity/api/v1/r/save", handleWithLog(capacity.SaveAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/list", handleWithLog(capacity.ListAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/get", handleWithLog(capacity.GetAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/delete", handleWithLog(capacity.DeleteAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/excel", handleWithLog(capacity.ExcelUploadHandler))
	http.Handle("/capacity/api/v1/plugin/compare", handleWithLog(capacity.ComparePluginHandler))
	http.Handle("/capacity/api/v1/export/expr/list", handleWithLog(capacity.ExportExprResult))
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/capacity/", http.StripPrefix("/capacity/", fs))
	listenPort := fmt.Sprintf(":%d", port)
	log.Logger.Info(fmt.Sprintf("listening %s ...", listenPort))
	http.ListenAndServe(listenPort, nil)
}

func handleWithLog(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		coreToken := models.CoreJwtToken{User: ""}
		requestLegal := true
		if models.Config().Http.AuthDisable != "true" {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				requestLegal = false
			} else {
				coreToken, requestLegal = authToken(tokenString)
				if strings.Contains(r.RequestURI, "/plugin/") {
					isSystemCall := false
					for _, v := range coreToken.Roles {
						if v == models.SystemRole {
							if coreToken.User == models.PlatformUser {
								isSystemCall = true
							}
							break
						}
					}
					requestLegal = isSystemCall
				}
			}
		}
		if requestLegal {
			h.ServeHTTP(w, r)
		} else {
			if strings.Contains(r.RequestURI, "/plugin/") {
				capacity.ReturnPluginAuthFail(r, w)
			} else {
				capacity.ReturnAuthFail(r, w)
			}
		}
		log.AccessLogger.Info("Request", log.String("url", r.RequestURI), log.String("method", r.Method), log.String("ip", strings.Split(r.RemoteAddr, ":")[0]), log.String("operator", coreToken.User), log.Float64("cost_second", time.Now().Sub(start).Seconds()))
	})
}

func authToken(token string) (result models.CoreJwtToken, legal bool) {
	key := models.CoreJwtKey
	result.User = ""
	legal = false
	if strings.HasPrefix(token, "Bearer") {
		token = token[7:]
	}
	if key == "" || strings.HasPrefix(key, "{{") {
		key = "Platform+Auth+Server+Secret"
	}
	keyBytes, err := ioutil.ReadAll(base64.NewDecoder(base64.RawStdEncoding, bytes.NewBufferString(key)))
	if err != nil {
		log.Logger.Error("Decode core token fail,base64 decode error", log.Error(err))
		return result, legal
	}
	pToken, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return keyBytes, nil
	})
	if err != nil {
		log.Logger.Error("Decode core token fail,jwt parse error", log.Error(err))
		return result, legal
	}
	claimMap, ok := pToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Logger.Error("Decode core token fail,claims to map error", log.Error(err))
		return result, legal
	}
	result.User = fmt.Sprintf("%s", claimMap["sub"])
	result.Expire, err = strconv.ParseInt(fmt.Sprintf("%.0f", claimMap["exp"]), 10, 64)
	if err != nil {
		log.Logger.Error("Decode core token fail,parse expire to int64 error", log.Error(err))
		return result, legal
	}
	roleListString := fmt.Sprintf("%s", claimMap["authority"])
	roleListString = roleListString[1 : len(roleListString)-1]
	result.Roles = strings.Split(roleListString, ",")
	return result, true
}
