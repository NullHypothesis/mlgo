include $(GOROOT)/src/Make.inc

TARG=mlgo/mlgo
GOFILES=\
				types.go\
				equal.go\
				stats.go\

include $(GOROOT)/src/Make.pkg
