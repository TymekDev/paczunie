window.onload = function() {
  document.querySelectorAll('.status-icons > .si').forEach(addIconClickHandler);
}

const addIconClickHandler = function(icon) {
  const id = icon.closest('li[id]').id;
  const s = getIconStatus(icon);
  icon.addEventListener('click', () => updateStatus(id, s));
}

const getIconStatus = function(icon) {
  return icon.className.match(/si-([1-3])/)[1];
}

const updateStatus = function(id, s) {
  const xhr = new XMLHttpRequest();
  xhr.addEventListener('load', () => {
    if (xhr.status !== 200) {
      // TODO: show a message to user that update failed
      return;
    }
    updatePkgStatusClassById(id, xhr.response);
  });
  xhr.open('PATCH', '/');
  xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
  xhr.send(`id=${id}&status=${s}`);
}

const updatePkgStatusClassById = function(id, s) {
  const pkg = document.getElementById(id);
  pkg.className = pkg.className.replace(/status-[1-3]/, `status-${s}`);
}
