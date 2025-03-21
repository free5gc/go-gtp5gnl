TOOL = gogtp5g-tunnel

.PHONY: all clean $(TOOL)

all: $(TOOL)

$(TOOL):
	CGO_ENABLED=0 go build -o $(TOOL) .

clean:
	go clean
	rm -f $(TOOL)
