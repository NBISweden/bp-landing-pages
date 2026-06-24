function toggleStainings(toggleElement) {
  const cell = toggleElement.parentElement;
  const hiddenItems = cell.querySelectorAll('.staining-item-hidden');
  const isExpanded = toggleElement.getAttribute('data-expanded') === 'true';

  hiddenItems.forEach((item) => {
    item.style.display = isExpanded ? 'none' : 'block';
  });

  toggleElement.textContent = isExpanded ? 'See more' : 'See less';
  toggleElement.setAttribute('data-expanded', isExpanded ? 'false' : 'true');
}

document.addEventListener('DOMContentLoaded', () => {
  const imageList = document.getElementById('dataset-example-images');

  if (imageList && typeof Viewer !== 'undefined') {
    new Viewer(imageList, { url: 'data-original' });
  }
});
