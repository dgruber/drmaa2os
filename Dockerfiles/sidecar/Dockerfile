FROM alpine:latest
WORKDIR /

COPY bin/* /

# https://kubernetes.io/docs/concepts/policy/pod-security-policy/#users-and-groups
USER 65534:65534

# sidecar waits for end of main container and does do some work
# (like epilogs or staging out files)
CMD [ "/sidecar" ]
