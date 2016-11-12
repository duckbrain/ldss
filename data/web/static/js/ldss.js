'use strict';

function log(e) { console.log(e); return e; }

var contentEle = document.querySelector('.main-content');
var previousBtn = document.getElementById('previous');
var nextBtn = document.getElementById('next');
var breadcrumbs = document.querySelector('.breadcrumbs');
var footnotes = document.querySelector('.footnotes');
var state = { item: null };

function interceptClickEvent(e) {
    var href;
    var target = e.target || e.srcElement;
    if (target.tagName === 'A' && !target.attributes.getNamedItem('disabled')) {
		e.preventDefault();
		
		href = target.getAttribute('href');
		if (href.indexOf('f_') == 0) {
			setFootnotesOpen(true);
			console.log("Footnote: " + href.substring(2));
		} else {
	        loadItem(target.pathname)
	        .then(function(item) {
				setFootnotesOpen(false);
				history.pushState(state, '', item.path);
			});
		}
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

function setFootnotesOpen(b) {
	document.body.classList.toggle('show-footnotes', b);
}

function setItem(item) {
	setNavigationButton(previousBtn, item.previous);
	setNavigationButton(nextBtn, item.next);
	
	breadcrumbs.innerHTML = '';
	item.breadcrumbs.forEach(function(crumb) {
		var content, a = document.createElement('a');
		a.classList.add('button');
		setNavigationButton(a, crumb);
		if (crumb.path == '/') {
			content = document.createElement('img');
			content.src = '/svg/home.svg';
			content.alt = "Library";
		} else {
			content = document.createTextNode(crumb.name)
		}
		a.appendChild(content);
		breadcrumbs.appendChild(a);
	})

	if (item.content) {
		if (item.content.indexOf('</h1>') != -1) {
			contentEle.innerHTML = item.content;
		} else {
			contentEle.innerHTML = "<h1>" + item.name + "</h1>" + item.content;
		}
	} else {
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
	
	footnotes.innerHTML = '';
	if (item.footnotes) {
		item.footnotes.forEach(function(footnote) {
			var li = document.createElement('li');
			li.textContent = footnote.name + ' - ' + footnote.linkName;
			li.innerHTML += '<div>' + footnote.content + '</div>';
			footnotes.appendChild(li);
		});
	}	
	
	state.item = item;
	return item;
}


document.addEventListener('click', interceptClickEvent);
window.addEventListener('popstate', onStateChange);