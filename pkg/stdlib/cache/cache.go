package cache

type CacheControl struct {
	MustRevalidate    bool
	SkipStoreNewCache bool
}

func NewCacheControl(cc *string) CacheControl {
	// when true then get from original source
	// when false then get from cache
	mustRevalidate := true
	if cc != nil {
		if *cc == "None" {
			mustRevalidate = false
		}
	}

	return CacheControl{
		MustRevalidate: mustRevalidate,
	}
}

type CacheControlBuilder interface {
	Build() CacheControl
	DefaultMustRevalidate(bool) CacheControlBuilder
	Parse(cc *string) CacheControlBuilder
}

type cacheControlBuilder struct {
	MustRevalidate bool
}

func NewCacheControlBuilder() CacheControlBuilder {
	return &cacheControlBuilder{}
}

func (cb *cacheControlBuilder) DefaultMustRevalidate(mustRevalidate bool) CacheControlBuilder {
	cb.MustRevalidate = mustRevalidate
	return cb
}

func (cb *cacheControlBuilder) Parse(cacheControl *string) CacheControlBuilder {
	if cacheControl == nil {
		// use default
		return cb
	}

	if *cacheControl == "Must Revalidate" {
		cb.MustRevalidate = true
	} else {
		cb.MustRevalidate = false
	}

	return cb
}

func (cb *cacheControlBuilder) Build() CacheControl {
	return CacheControl{
		MustRevalidate: cb.MustRevalidate,
	}
}
