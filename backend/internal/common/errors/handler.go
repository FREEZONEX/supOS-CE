package errors

//// GlobalErrorHandler returns a middleware that recovers from any panics and handles errors.
//// This middleware is designed to work with the standard net/http library, making it compatible with frameworks like go-zero.
//func GlobalErrorHandler(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		defer func() {
//			if err := recover(); err != nil {
//				// Log panic with stack trace
//				log.Printf("panic recovered: %v\n%s", err, debug.Stack())
//
//				// First, assert that the recovered value is an error.
//				if recoveredErr, ok := err.(error); ok {
//					// Then, check if it's our specific business error.
//					if buzErr, isBuz := IsBuzError(recoveredErr); isBuz {
//						httpStatus := http.StatusBadRequest // Default for business errors
//						if buzErr.Code >= 400 && buzErr.Code < 600 {
//							httpStatus = buzErr.Code
//						}
//						msg := I18nUtils.GetMessage(buzErr.Msg, buzErr.Params...)
//						Fail(w, httpStatus, buzErr.Code, msg)
//						return
//					}
//				}
//
//				// Handle unknown errors (non-BuzError panics or generic errors)
//				Fail(w, http.StatusInternalServerError, http.StatusInternalServerError, "Internal Server Error")
//			}
//		}()
//		// Call the next handler in the chain
//		next.ServeHTTP(w, r)
//	})
//}
