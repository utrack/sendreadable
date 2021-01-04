package htmllatex

import "context"

func isOuterPar(ctx context.Context) bool {
	v, ok := ctx.Value(ctxKeyNotOuterPar).(bool)
	if !ok {
		return true
	}
	return !v
}

func isParable(ctx context.Context) bool {
	v, ok := ctx.Value(ctxKeyNotParable).(bool)
	if !ok {
		return true
	}
	return !v
}

func cntNotParable(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyNotParable, true)
}
func cntNotOuterPar(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyNotOuterPar, true)
}

var ctxKeyNotOuterPar = "notOuterPar"
var ctxKeyNotParable = "notParable"
