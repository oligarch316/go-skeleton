package command

// Option TODO
type Option func(*Options)

// WithAppendFields TODO
func WithAppendFields(fields ...Field) Option {
	return func(o *Options) { o.Fields = append(o.Fields, fields...) }
}

// WithReplaceFields TODO
func WithReplaceFields(fields ...Field) Option {
	return func(o *Options) { o.Fields = fields }
}
