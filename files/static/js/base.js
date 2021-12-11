
window.al = function(x) {
    alert(x);
}

window.swapVisibility = function (el) {
    // let style = window.getComputedStyle(el);
    if( el.style.getPropertyValue('visibility') === 'visible' ) {
        el.style.setProperty('visibility', 'hidden');
    } else {
        el.style.setProperty('visibility', 'visible');
    }
}

window.hamburger = function (open, closed, content ) {
    let content_el = document.getElementById(content);
    if( content_el.style.getPropertyValue('visibility') !== 'visible' ) {
        document.body.scrollTop = 0; // For Safari
        document.documentElement.scrollTop = 0; // For Chrome, Firefox, IE and Opera
    }
    window.swapVisibility(document.getElementById(open));
    window.swapVisibility(document.getElementById(closed));
    window.swapVisibility(content_el);

}


function firstElemChild(el) {
    if (el.firstElementChild)
        return el.firstElementChild;

    el = el.firstChild

    while (el && el.nodeType !== 1)
        el = el.nextSibling;

    return el;
}

/**
 * @preserve  textfill
 * @name      jquery.textfill.js
 * @author    Russ Painter (GeekyMonkey)
 * @author    Yu-Jie Lin
 * @author    Alexandre Dantas
 * @version   0.6.2
 * @date      2018-02-24
 * @copyright (c) 2009
 * @license   MIT License
 * @homepage  https://github.com/jquery-textfill/jquery-textfill
 * @example   http://jquery-textfill.github.io/jquery-textfill/index.html
 */

function textHeight(el) {
    let height = Math.ceil(el.offsetHeight);
    return height;
}

function textWidth(el) {
    let width = Math.ceil(el.offsetWidth);
    return width;
}


/**
 * Calculates which size the font can get resized,
 * according to constrains.
 *
 * @param {String} prefix Gets shown on the console before
 *                        all the arguments, if debug mode is on.
 * @param {Object} ourText The DOM element to resize,
 *                         that contains the text.
 * @param {function} func Function called on `ourText` that's
 *                        used to compare with `max`.
 * @param {number} max Maximum value, that gets compared with
 *                     `func` called on `ourText`.
 * @param {number} minFontPixels Minimum value the font can
 *                               get resized to (in pixels).
 * @param {number} maxFontPixels Maximum value the font can
 *                               get resized to (in pixels).
 *
 * @return Size (in pixels) that the font can be resized.
 */
function _sizing(prefix, ourText, func, max, maxHeight, maxWidth, minFontPixels, maxFontPixels) {

    // The kernel of the whole plugin, take most attention
    // on this part.
    //
    // This is a loop that keeps increasing the `font-size`
    // until it fits the parent element.
    //
    // - Start from the minimal allowed value (`minFontPixels`)
    // - Guesses an average font size (in pixels) for the font,
    // - Resizes the text and sees if its size is within the
    //   boundaries (`minFontPixels` and `maxFontPixels`).
    //   - If so, keep guessing until we break.
    //   - If not, return the last calculated size.
    //
    // I understand this is not optimized and we should
    // consider implementing something akin to
    // Daniel Hoffmann's answer here:
    //
    //     http://stackoverflow.com/a/17433451/1094964
    //


    ourText.style.removeProperty("font-size");
    let styleval = window.getComputedStyle(el, null).getPropertyValue("font-size");
    let fontSize = parseFloat(styleval);
    //let fontSize = ourText.style.fontSize;
    let lastFontSize = fontSize;
    let counter = 0;
    ourText.style.setProperty('font-size', styleval)
    styleval = window.getComputedStyle(el, null).getPropertyValue("font-family");
    ourText.style.setProperty('font-family', styleval)
    let width = ourText.offsetWidth;
    let height = ourText.offsetHeight;


    while (width < maxWidth-1 && height < maxHeight-1 && counter < 200) {
        counter++;
        lastFontSize = fontSize;
        fontSize+=2;
        //ourText.style.removeProperty("font-size");
        ourText.style.setProperty("font-size", fontSize+'px');
        width = ourText.offsetWidth;
        height = ourText.offsetHeight;
    }
    return lastFontSize;
}
function _sizing2(prefix, ourText, func, max, maxHeight, maxWidth, minFontPixels, maxFontPixels) {
    ourText.style.removeProperty("font-size");
    let styleval = window.getComputedStyle(el, null).getPropertyValue("font-size");
    let fontSize = parseFloat(styleval);
    //let fontSize = ourText.style.fontSize;
    let lastFontSize = fontSize;
    let counter = 0;
    ourText.style.setProperty('font-size', styleval)
    styleval = window.getComputedStyle(el, null).getPropertyValue("font-family");
    ourText.style.setProperty('font-family', styleval)
    let width = ourText.offsetWidth;
    let height = ourText.offsetHeight;

    let text = ourText.innerText;

    while (width <= maxWidth+1 && height <= maxHeight+1 && counter < 200) {
        counter++;
        lastFontSize = fontSize;
        fontSize+=2;
        //ourText.style.removeProperty("font-size");
        ourText.style.setProperty("font-size", fontSize+'px');
        ourText.innerHTML = text;
        width = ourText.offsetWidth;
        height = ourText.offsetHeight;
    }
    ourText.style.setProperty("font-size", lastFontSize+'px');
    return lastFontSize;
}

window.optimizeFontsize = function () {

    let testbox = document.getElementById("testbox");
    let testfield = document.getElementById("testfield");

    let els = document.getElementsByClassName("type-cell");
    let els2 = document.getElementsByClassName("type-news");
    els = [...els, ...els2];
    for ( i = 0; i < els.length; i++ ) {
        el = els[i];
        let daText = el.innerText;
        if ( daText == '' ) { continue; }

        let maxSize = 0;
        let maxHeight = el.clientHeight;
        let maxWidth  = el.clientWidth;
        testfield.innerText = daText;

        //let rect = el.getBoundingClientRect()
        testbox.style.setProperty('width', maxHeight+'px' );
        testbox.style.setProperty('height', maxWidth+'px' );

        let styleval = window.getComputedStyle(el, null).getPropertyValue('font-size');
        let oldFontSize = parseFloat(styleval);

        let fontSizeWidth = undefined;

        let text = el.querySelector('.text');
        fontSizeWidth = _sizing2(
            'Width', text,
            textWidth, maxWidth,
            maxHeight, maxWidth,
            8, maxHeight
        );

        if ( fontSizeWidth > maxSize ) {
            maxSize = fontSizeWidth;
        }
    }
}

window.decorateText = function(el) {
    let t = document.getElementById("content-text");
    // todo: create some nice logic...
    let range = document.createRange();
    range.setStart(myElementWithText.childNodes[myTextElementIndex], myLetterPosition);
    range.setEnd(myElementWithText.childNodes[myTextElementIndex], myLetterPosition+1);
    return range.getBoundingClientRect();

}