function initializeTabs(element) {
  function childElements(element) {
    return [].slice.call(element.childNodes)
      .filter(function(el) {
        return el.nodeType === 1;
      });
  }

  var tabs = childElements(element.querySelector('ul.tabs'));
  var tabContents = childElements(element.querySelector('span.tabs-body'));

  tabs.forEach(function(tabElement, i) {
    tabElement.classList.add('highlight');
    tabElement.dataset.id = i;

    tabElement.addEventListener('click', function() {
      tabs.forEach(function(tabElement) {
        tabElement.classList.remove('active');
      });

      tabContents.forEach(function(tabContent) {
        tabContent.classList.remove('active');
        tabContent.classList.add('hide');

        if (tabElement.dataset.id === tabContent.dataset.id) {
          tabContent.classList.add('active');
          tabContent.classList.remove('hide');
        }
      });

      tabElement.classList.add('active');
    });
  });

  tabContents.forEach(function(tabContent, i) {
    tabContent.classList.add('highlight', 'hide');
    tabContent.dataset.id = i;
  });

  tabs[0].classList.add('active');
  tabContents[0].classList.add('active');
  tabContents[0].classList.remove('hide');
}
