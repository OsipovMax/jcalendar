package jcalendar

func pcaster[T uint | string](val T) *T {
	return &val
}
