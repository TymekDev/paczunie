window.onload = function() {
  document.querySelectorAll('.status-icons > .si').forEach(addIconClickHandler);
}

const addIconClickHandler = function(icon) {
  const pkg = icon.closest('.package');
  const id = pkg.id;
  const s = getIconStatus(icon);
  icon.addEventListener('click', () => {
    if (getPkgStatus(pkg) !== s) {
      updateStatus(id, s);
    }
  });
}

const getIconStatus = function(icon) {
  return icon.className.match(/si-([1-9])/)[1];
}

const getPkgStatus = function(pkg) {
  return pkg.className.match(/status-([1-9])/)[1];
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
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.send(`id=${id}&status=${s}`);
}

const updatePkgStatusClassById = function(id, s) {
  const pkg = document.getElementById(id);
  pkg.className = pkg.className.replace(/status-[1-9]/, `status-${s}`);
}

const deletePkg = function(id) {
  const xhr = new XMLHttpRequest();
  xhr.addEventListener('load', () => {
    if (xhr.status !== 200) {
      // TODO: show a message to user that delete failed
      return;
    }
    document.getElementById(id).remove();
  });
  xhr.open('DELETE', '/');
  xhr.send(`${id}`);
}
