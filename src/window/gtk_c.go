package window
/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>

extern void goHandleEvent(void * data);

GdkEventType GetGTKEventType(GdkEvent * event) {
	return event->type;
}

void GtkSetEvent(gpointer * instance, char * event, void * data) {
	g_signal_connect_swapped(instance, event, goHandleEvent, data);
}
*/
import "C"