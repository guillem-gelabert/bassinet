# Bassinet

Bassinet is a set of 11 utility middlewares to help secure HTTP headers. It's based on the widely used helmet.js. Includes middleware functions for setting the following headers:

- `X-XSS-Protection`
- `Strict-Transport-Security`
- `Referrer-Policy`
- `X-Permitted-Cross-Domain-Policies`
- `X-Download-Options`
- `X-Powered-By`
- `X-Frame-Options`
- `Expect-CT`
- `X-Content-Type-Options`
- `X-DNS-Prefetch-Control`
- `Content-Security-Policy`

## Usage

Initialize the middleware with the desired options —if any— and handle the returned error.

```
referrerPolicy, err := bassinet.ReferrerPolicy([]{
	bassinet.PolicyOrigin,
	bassinet.PolicyUnsafeURL,
})
if err != nil {
    // handle error
}
```

### With ServeMux

To use bassinet with the builtin ServeMux you just wrap it with the initialized middleware.

```
mux := http.NewServeMux()
mux.HandleFunc("/", home)
srv := http.Server{
	Handler: referrerPolicy(mux)
}
```

As you might probably want to chain several of the middlewares it is recommended to use a composing function.

###Usage with justinas/alice

```
xssFilter, err := bassinet.XSSFilter()
if err != nil {
    // handle error
}
htsts, err := bassinet.StrictTransportSecurity(StrictTransportOption{
	maxAge:            60,
  excludeSubdomains: true,
})
if err != nil {
	// handle error
}

middleware := alice.New(xssFilter, htsts)
```



