'use strict';

function log(e) { console.log(e); return e; }

var contentEle = document.querySelector('.main-content');
var toolbar = document.querySelector('.toolbar');
var previousBtn = document.getElementById('previous');
var nextBtn = document.getElementById('next');
var breadcrumbs = document.querySelector('.breadcrumbs');
var footnotes = document.querySelector('.footnotes');
var footnotesHeader = document.querySelector('.footnotes-header');
document.querySelector('.footnotes-close').addEventListener('click', function(e) {
	e.preventDefault();
	setFootnotesOpen(false);
});
var state = { item: null };

var scrolledVerse = null;
var scrolledVerseNumber = 0;
var lastScrollPosition = 0;
var scrollTimeout = null;

function interceptClickEvent(e) {
    var href;
    var target = e.target || e.srcElement;
    if (target.tagName !== 'A' || 
		target.attributes.getNamedItem('disabled') ||
		e.ctrlKey || e.button !== 0){
			return;
	}
	
	href = target.getAttribute('href');
	if (href.indexOf('f_') == 0) {
		e.preventDefault();
		setFootnotesOpen(true);
		setFootnote(href.substring(2));
		console.log("Footnote: " + href.substring(2));
	} else {
		return;
        loadItem(target.pathname)
        .then(function(item) {
			setFootnotesOpen(false);
			history.pushState(state, '', item.path);
		});
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

function setFootnote(ref) {
	var element = document.getElementById('ref-' + ref);
	if (element == null) {
		return;
	}
	scrollTo(footnotes, element.offsetTop - footnotesHeader.clientHeight, 200);
}

function scrollTo(element, to, duration) {
    if (duration <= 0) return;
    var difference = to - element.scrollTop;
    var perTick = difference / duration * 10;

    setTimeout(function() {
        element.scrollTop = element.scrollTop + perTick;
        if (element.scrollTop === to) return;
        scrollTo(element, to, duration - 10);
    }, 10);
}

function scrollToHighlight() {
	var el = document.querySelector('.highlight');
	if (!el || location.hash) {
		return;
	}
	scroll(0, el.offsetTop - footnotesHeader.clientHeight);
}

function onScroll() {
	var s = window.scrollY;
	if (lastScrollPosition == s || (scrolledVerse != null && isScrolledOnto(scrolledVerse))) {
		return
	}
	
	if (lastScrollPosition < s) {
		do {
			scrolledVerseNumber++;
			scrolledVerse = document.getElementById(scrolledVerseNumber);
			if (!scrolledVerse) {
				scrolledVerseNumber = 0;
				lastScrollPosition = 0;
				return
			}
		} while (!isScrolledOnto(scrolledVerse));
		lastScrollPosition = s;
	} else {
		do {
			scrolledVerseNumber--;
			scrolledVerse = document.getElementById(scrolledVerseNumber);
			if (!scrolledVerse) {
				scrolledVerseNumber = 0;
				lastScrollPosition = 0;
				return
			}
		} while (!isScrolledOnto(scrolledVerse));
		
		lastScrollPosition = s;
	}
	setFootnote(scrolledVerseNumber + 'a');
}

function isScrolledOnto(elem) {	
	var y = elem.offsetTop;
    var height = elem.offsetHeight;

    while ( elem = elem.offsetParent )
        y += elem.offsetTop;
	y -= toolbar.clientHeight;

    var maxHeight = y + height;
    var isVisible = ( y < ( window.pageYOffset) ) && ( maxHeight >= window.pageYOffset );
    return isVisible; 
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
		item.footnotes.forEach(function(footnote, i) {
			var li = document.createElement('li');
			li.textContent = footnote.name + ' - ' + footnote.linkName;
			if (i > 0 && footnote.name.indexOf(item.footnotes[i-1].name + " ") == 0) {
				footnote.name = footnote.name.substring(item.footnotes[i-1].name.length + 1);
			}
			li.innerHTML += '<div>' + footnote.content + '</div>';
			footnotes.appendChild(li);
		});
	}	
	
	state.item = item;
	return item;
}

document.addEventListener('click', interceptClickEvent);
scrollToHighlight();
window.addEventListener('popstate', onStateChange);
window.addEventListener('scroll', function() {
	clearTimeout(scrollTimeout);
	scrollTimeout = setTimeout(onScroll, 100);
});
