package lib

func Parse(lang *Language, q string) (r Reference, err error) {
	if ref, e := lang.ref(); e == nil {
		r, err = ref.lookup(q)
		r.Language = lang
		return
	}
	return ParsePath(lang, q)
}

func ParsePath(lang *Language, p string) (Reference, error) {
	return Reference{
		Language: lang,
		GlPath:   p,
	}, nil
}
