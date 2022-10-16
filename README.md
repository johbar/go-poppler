# go-poppler

A simple wrapper for the [poppler pdf lib](https://poppler.freedesktop.org/) for Golang.

This go module exports a limited subset of Popplers functionalties.

## Install dependencies

In order to use this module you need to install two packages. In Debian/Ubuntu this is done by:

```sh
apt-get install libpoppler-glib-dev libcairo2-dev
```

## Usage

```sh
go get github.com/johbar/go-poppler
```

## Improvements and performance considerations

This is a fork of [timsat/go-poppler](/timsat/go-poppler) which is derived of [cheggaaa/go-poppler](cheggaaa/go-poppler).

This fork uses [finalizers](https://pkg.go.dev/runtime#SetFinalizer) as a safeguard against memory leaks.
This enables you not to call `document.Close()` and `page.Close()` but let the GC do the clean-up work instead.

This might add a significant memory overhead in long-running high-throughput processes (like a webservice or batch processor), when a lot of *Poppler* objects in unmanaged/off-heap memory are being created. They are evicted by the finalizer but the GC might run too late to prevent a OOM as it doesn't take these into account when it schedules the next cycle.

One advantage of relying on finalizes solely might be an improved CPU utilization because they run on their own goroutine. So your main routine doesn't need to handle the clean-up. (But you can archive this in many other ways, I guess, e.g. by writing your own goroutine and using channels.)

*Tl;dr*: Be careful when relying on finalizers; do some (load) tests and watch your RAM!
