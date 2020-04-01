FROM scratch
COPY hfs /
VOLUME ["/htdocs"]
ENTRYPOINT ["/hfs"]
CMD ["-addr", ":8000", "-dotfile", "-root", "/htdocs"]
