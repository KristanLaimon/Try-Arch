# Runtime validation

Before changing the pinned QEMU runtime, run runtime-build validation tests and
check Windows CPU boot, accelerated boot, keyboard focus, audio routing, USB,
display resize, shutdown, reboot and suspend/resume. The existing portable
runtime has a frozen environment-variable API; keep those names compatible
until the binary and launcher are released together. These names are not guest
repositories or desktop customizations.
