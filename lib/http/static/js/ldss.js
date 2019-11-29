'use strict';

function log(e) { console.log(e); return e; }

var contentEle = document.querySelector('.main-content');
var toolbar = document.querySelector('.toolbar');
var previousBtn = document.getElementById('previous');
var nextBtn = document.getElementById('next');
var breadcrumbs = document.querySelector('.breadcrumbs');
var footnotes = document.querySelector('.footnotes');
var footnotesHeader = document.querySelector('.footnotes-header');
var searchBox = document.querySelector("[name=q]")

document.querySelector('.footnotes-close').addEventListener('click', function(e) {
	e.preventDefault();
	setFootnotesOpen(false);
});
var state = { item: null };

var scrolledVerse = null;
var scrolledVerseNumber = 0;
var lastScrollPosition = 0;
var scrollTimeout = null;

var verseSearch = '';


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

function scrollToTop(element, val) {
	if (element === window) {
		if (val) {
			scroll(0, val);
		} else {
			return element.scrollY;
		}
	} else {
		if (val) {
			element.scrollTop = val
		} else {
			return element.scrollTop;
		}
	}
}


let scrollToTimer = null
function scrollTo(element, to, duration) {
    if (duration <= 0) return;
    var difference = to - scrollToTop(element);
    var perTick = difference / duration * 10;

    clearTimeout(scrollToTimer);
    scrollToTimer = setTimeout(function() {
        scrollToTop(element, scrollToTop(element) + perTick);
        if (scrollToTop(element) === to) return;
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

function scrollVerse(n) {
	const elsVerses = document.querySelectorAll('.verse')
	const scrollY = window.pageYOffset;
	const buffer = 10;
	let verse = null;
	let verse1 = document.getElementById(1);
	if (verse1 && scrollY < verse1.offsetTop) {
		verse = 0;
	} else {
		for (let el of elsVerses) {
			if (el && scrollY >= el.offsetTop) {
				//&& scrollY <= (el.offsetTop + el.clientHeight + buffer)) {
				verse = parseInt(el.id);
			} else {
				break;
			}
		}
	}

	if (verse === null) {
		return;
	}

	if (verse + n === 0) {
		scrollTo(window, 0, 200);
	}
	const elVerse = document.getElementById(verse + n)
	if (elVerse) {
		scrollTo(window, elVerse.offsetTop , 200);
	} else if (n > 0) {
		scrollTo(window, document.clientHeight, 200);
	}
}

function onKeyPress(e) {
	if (e.target === searchBox) {
		return;
	}

	let el;
	let handled = true
	if (!e.altKey && !e.ctlKey && !e.shiftKey) {
		switch (e.keyCode) {
			case 8: //backspace
				el = document.querySelector(".breadcrumbs .button:nth-last-child(2)")
				if (el) el.click()
				break;
			case 39: // right-arrow
			case 76: // 'l' key
				scrollVerse(1);
				break;
			case 37: // left-arrow
			case 72: // 'h' key
				scrollVerse(-1);
				break;
			case 74: // 'j' key
				scrollTo(window, window.pageYOffset + 100, 70);
				break;
			case 75: // 'k' key
				scrollTo(window, window.pageYOffset - 100, 70);
				break;
			case 78: // 'n' key
				document.getElementById('next').click();
				break;
			case 66: // 'b' key
				document.getElementById('previous').click();
				break;

			case 48:
			case 49:
			case 50:
			case 51:
			case 52:
			case 53:
			case 54:
			case 55:
			case 56:
			case 57:
				verseSearch += (e.keyCode - 48);
				while (!(el = document.getElementById(verseSearch)) && verseSearch != '') {
					verseSearch = verseSearch.substring(1)
				}
				console.log(verseSearch);
				const verse = parseInt(verseSearch)
				if (Number.isNaN(verse) || !verse) {
					scrollTo(window, 0, 200);
				}
				if (el) {
					scrollTo(window, el.offsetTop, 200)
				}
				break;
				


			case 191: // '/' key
				el = document.querySelector('.lookup input');
				el.focus()
				el.select()
				e.preventDefault()
				break;
			default: handled = false; break;
		}
	} else {
		handled = false;
	}

	if (handled) {
		e.preventDefault();
	} else {
		console.log(e.keyCode, e)
	}
}

document.addEventListener('click', interceptClickEvent);
scrollToHighlight();
window.addEventListener('popstate', onStateChange);
window.addEventListener('scroll', function() {
	clearTimeout(scrollTimeout);
	scrollTimeout = setTimeout(onScroll, 100);
});
window.addEventListener('keydown', onKeyPress);
