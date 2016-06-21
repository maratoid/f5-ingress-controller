FROM gcr.io/google_containers/busybox
COPY controller /
CMD ["/controller"]