package transport

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

//func (h *Handler) userIdentity(c *gin.Context) {
//	header := c.GetHeader(authorizationHeader)
//	if header == "" {
//		newErrorResponse(c, http.StatusUnauthorized, "empty authorization header")
//		return
//	}
//	headerParts := strings.Split(header, " ")
//	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
//		newErrorResponse(c, http.StatusUnauthorized, "invalid authorization header")
//		return
//	}
//	if len(headerParts[1]) == 0 {
//		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
//		return
//	}
//	// parse token
//	userId, err := h.service.ParseToken(headerParts[1])
//	if err != nil {
//		newErrorResponse(c, http.StatusUnauthorized, err.Error())
//		return
//	}
//	// Write value id in context (need for access id in another service (for another  handlers))
//	c.Set(userCtx, userId)
//}
//
//func getUserId(c *gin.Context) (int, error) {
//	userId, ok := c.Get(userCtx)
//	if !ok {
//		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
//		return 0, errors.New("user id not found")
//	}
//	idInt, ok := userId.(int)
//	if !ok {
//		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
//		return 0, errors.New("user id is of invalid type")
//	}
//	return idInt, nil
//}
