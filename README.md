# Bassinet

Bassinet is a set of 11 utility middlewares to help secure HTTP headers. It's based on the widely used helmet.js. Includes middleware functions for setting the following headers:

- [`X-XSS-Protection`](#X-XSS-Protection)
- [`Strict-Transport-Security`](#Strict-Transport-Security)
- [`Referrer-Policy`](#Referrer-Policy)
- [`X-Permitted-Cross-Domain-Policies`](#X-Permitted-Cross-Domain-Policies)
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

### Response Headers

### X-XSS-Protection

`XSSFilter` sets `X-XSS-Protection` header to `0` to prevent attackers from blocking legit code or inferring information. [Read more](https://guillem-gelabert.github.io/posts/x-xss-protection/). `XSSFilter` accepts no options.

```
xssFilter, err := bassinet.XSSFilter()
if err != nil {
	// Handle error
}

srv := http.Server{
	Handler: xssFilter(mux)
}
```

### Strict-Transport-Security

StrictTransportSecurity sets Strict-Transport-Security so that browsers remember if HTTPS is available, to avoid insecure connection before redirect. [Read more](https://guillem-gelabert.github.io/posts/strict-transport-security/).

It accepts a `bassinet.StrictTransportOptions` struct to set the following directives:

- `maxAge`: Time (in seconds) that the browser should remember if the site has HTTPS. Defaults to 180 days. **int**
- `excludeSubdomains`: Optional. If set the browser will apply directive to subdomains. **bool**
- `preload`: Optional. If set the browser will check the [Preloading Strict Transport Security](https://hstspreload.org/) public list, enabling STS also on first load. **bool**

```
policies := bassinet.StrictTransportOptions{
	maxAge:            60 * 60 * 24 * 7, // recheck every week
	excludeSubdomains: true,
	preload: true,
}

sts, err := bassinet.StrictTransportSecurity(policies)
if err != nil {
	// Handle error
}

srv := http.Server{
	Handler: sts(mux)
}
```

### X-Permitted-Cross-Domain-Policies

PermittedCrossDomainPolicies sets X-Permitted-Cross-Domain-Policies header to tell some user-agents (most notably Adobe products) your domain's policy for loading cross-domain content. [Read more](https://www.adobe.com/devnet-docs/acrobatetk/tools/AppSec/xdomain.html).

Accepts the following policies:

- `PCDPNone`: No `crossdomain.xml` file is allowed.
- `PCDPMasterOnly`: Only check `crossdomain.xml` in the root directory of the website.
- `PCDPByContentType`: Only accept files with type `text/x-cross-domain-policy`.
- `PCDPAll`: Allow any `crossdomain.xml` files.

```
permittedCrossDomainPolicies, err := bassinet.PermittedCrossDomainPolicies(bassinet.PCDPByContentType)
if err != nil {
	// Handle error
}

srv := http.Server{
	Handler: permittedCrossDomainPolicies(mux)
}
```
