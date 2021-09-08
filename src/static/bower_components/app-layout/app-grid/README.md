##&lt;app-grid&gt;

app-grid is a helper class useful for creating responsive, fluid grid layouts using custom properties.
Because custom properties can be defined inside a `@media` rule, you can customize the grid layout
for different breakpoints.

Example:

Within an element definition, import `app-grid-style.html`, then include the style sheet in the style
element inside the template.

```html
<template>
  <style include="app-grid-styles">

    :host {
      --app-grid-columns: 3;
      --app-grid-item-height: 100px;
    }

    .item {
      background-color: white;
    }

    @media (max-width: 640px) {
      :host {
        --app-grid-columns: 1;
      }
    }

  </style>

  <div class="app-grid">
    <div class="item">1</div>
    <div class="item">2</div>
    <div class="item">3</div>
  </div>
</template>
```

In this example, the grid  will take 3 columns per row and only 1 column if the viewport width is
smaller than 640px.

### Expandible items

In many cases, it's useful to expand an item more than 1 column. To achieve this type of layout,
you can specify the number of columns the item should expand to by setting the custom property
`--app-grid-expandible-item`. To indicate which item should expand, apply the mixin
`--app-grid-expandible-item` to a rule with a selector to the item. For example:

```css
<style>

  :host {
    --app-grid-columns: 3;
    --app-grid-item-height: 100px;
    --app-grid-expandible-item: 3;
  }

  /* Only the first item should expand */
  .item:first-child {
    @apply(--app-grid-expandible-item);
  }

</style>
```

### Preserving the aspect ratio

When the size of a grid item should preserve the aspect ratio, you can add the `has-aspect-ratio`
attribute to the element with the class `.app-grid`. Now, every item element becomes a wrapper around
the item content. For example:

```css
<template>
  <style>

    :host {
      --app-grid-columns: 3;
      /* 50% the width of the item is equivalent to 2:1 aspect ratio*/
      --app-grid-item-height: 50%;
    }

    .item {
      background-color: white;
    }

  </style>
  <div has-aspect-ratio>
    <div class="item">
      <div>item 1</div>
    </div>
    <div class="item">
      <div>item 2</div>
    </div>
    <div class="item">
      <div>item 3</div>
    </div>
  </div>
</template>
```

### Styling

Custom property                               | Description                                                | Default
----------------------------------------------|------------------------------------------------------------|------------------
`--app-grid-columns`                          | The number of columns per row.                             | 1
`--app-grid-gutter`                           | The space between two items.                               | 0px
`--app-grid-item-height`                      | The height of the items.                                   | auto
`--app-grid-expandible-item-columns`          | The number of columns an expandible item should expand to. | 1
