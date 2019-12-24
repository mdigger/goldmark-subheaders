# goldmark-subheaders

[Goldmark](https://github.com/yuin/goldmark/) extension for markdown sub-headers.

```markdown
## Title
## Subtitle
text

## Subtitle 1

## Subtitle 2

text
```

```html
<header>
<h2>Title</h2>
<p>Subtitle</p>
</header>
<p>text</p>
<h2>Subtitle 1</h2>
<h2>Subtitle 2</h2>
<p>text</p>
```

```go
// add sub-headers support to golmark markdown parser
md := goldmark.New(subheaders.Option)
```
