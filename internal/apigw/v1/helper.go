package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/status"

	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/httputil"
)

func handleGRPCError(w http.ResponseWriter, err error) {
	slog.Error("handleGRPCError", slog.String("err", err.Error()))

	w.Header().Set("Content-Type", "application/json")

	st := status.Convert(err)
	code := st.Code()
	w.WriteHeader(httputil.ConvertGRPCCodeToHTTP(code))
	if err := json.NewEncoder(w).Encode(
		apiv1.Error{
			Code:    httputil.ConvertGRPCToErrorCode(code),
			Message: nil,
		},
	); err != nil {
		slog.Error(" json.NewEncoder Encode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
