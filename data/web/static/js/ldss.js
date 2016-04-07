function log(e) { console.log(e); return e; }

var contentEle = document.querySelector('.main-content');
var previousBtn = document.getElementById('previous');
var nextBtn = document.getElementById('next');
var parentBtn = document.getElementById('parent');
var state = { item: null };

function interceptClickEvent(e) {
    var href;
    var target = e.target || e.srcElement;
    if (target.tagName === 'A' && !target.attributes.getNamedItem('disabled')) {
		href = target.getAttribute('href');
        loadItem(target.pathname)
        .then(function(item) {
			history.pushState(state, '', item.path);
		});

	   e.preventDefault();
    }
}

function onStateChange(e) {
	if (e.state)
		setState(e.state);
	else {
		loadItem(location.pathname);
	}
}

function setNavigationButton(element, item) {
	element.title = item ? "enabled" : "disabled";
	if (item) {
		element.removeAttribute('disabled');
		element.href = item.path;
	} else {
		element.removeAttribute('href');
		element.setAttribute('disabled', '')
	}
}

function loadItem(pathname) {
	return fetch("/api" + pathname).then(function(res) {
		if (res.status != 200) throw res;
		return res.json().then(setItem);
	});
}

function setState(state) {
	setItem(state.item);
}

function setItem(item) {
	setNavigationButton(previousBtn, item.previous);
	setNavigationButton(nextBtn, item.next);
	setNavigationButton(parentBtn, item.parent);

	if (item.content)
		contentEle.innerHTML = item.content;
	else {
		contentEle.innerHTML = '';
		var header = document.createElement('h1');
		var children = document.createElement('ul');

		header.classList.add('item-name');
		header.textContent = item.name;
		children.classList.add('item-children');

		item.children.forEach(function(child) {
			var childEle = document.createElement('li');
			var childLink = document.createElement('a');

			childLink.href = child.path+"?lang="+child.language;
			childLink.textContent = child.name;

			childEle.appendChild(childLink)
			children.appendChild(childEle)
		});

		contentEle.appendChild(header);
		contentEle.appendChild(children);

	}
	state.item = item;
	return item;
}


document.addEventListener('click', interceptClickEvent);
window.addEventListener('popstate', onStateChange);