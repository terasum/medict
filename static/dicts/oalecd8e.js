var oalecd8e = {};
oalecd8e.query = jQuery.noConflict( true );

const OxfordTagSwitchCN = '.oalecd8e_switch_lang';
const OxfordTagSwitchCNALL = '.oalecd8e_switch_lang.switch_all';
const OxfordTagSwitchCNCHILDREN = '.oalecd8e_switch_lang.switch_children';
const OxfordTagSwitchCNSIBLINGS = '.oalecd8e_switch_lang.switch_siblings';
const OxfordTagSwitchCNTAG = '.oalecd8e_switch_lang.switch_children, .oalecd8e_switch_lang.switch_siblings'

const OxfordTagChineseTexT = '.oalecd8e_chn';

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!! script embedded in HTML !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// add platfrom infos to <span collinsbody>

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!    embedd jQuery 3.2.1 !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
/*! jQuery v3.2.1 | (c) JS Foundation and other contributors | jquery.org/license */
/* beautify ignore:start */

/* beautify ignore:end */
// -------------------------------- END END END jQuery 3.2.1 -------------------------------

oalecd8e.query(oalecd8e_documentReady);

var oalecd8e_pagetype = 0;
var oalecd8e_pageCount = 0;

function oalecd8e_documentReady() {
  oalecd8e_extendJQuery();
  oalecd8e_switchChineseSetup();
  oalecd8e_dblSetup();

  oalecd8e.query(window).scroll(); // trigger float logo

}

var oalecd8e_lastdblSetupClick = null;

function oalecd8e_dblSetup() {
  oalecd8e.query('.entry').off('.entry').on('click.entry', function(event) {
    oalecd8e.query('.n-g, .x-g, .sense-g, .para, .def-g')
      .each(function() {
        // console.log('classname:' + this.className + " " + oalecd8e.query(this).offset().top + ' ' + (oalecd8e.query(this).offset().top + oalecd8e.query(this).height()) + " " + event.pageY);
        if ((oalecd8e.query(this).offset().top < event.pageY) &&
          ((oalecd8e.query(this).offset().top + oalecd8e.query(this).outerHeight()) > event.pageY)) {
          // console.log('target:' + oalecd8e.query(event.target).filter('.n-g').children(".def-g").first().children(OxfordTagSwitchCNTAG)
          //                                 .first()
          //                                 .trigger('entry').length);
          var _element = oalecd8e.query(this);
          if (_element.is('.n-g, .x-g, .sense-g, .para, .def-g')) {
            if (oalecd8e.query(event.target)
              .filter('.sense-g, .n-g')
              .children('.def-g')
              .first()
              .children(OxfordTagSwitchCNTAG)
              .first()
              .trigger('entry').length > 0
            ) return false;

            if (oalecd8e.query(event.target)
              .filter('.x-g, .para, .def-g')
              .first()
              .children(OxfordTagSwitchCNTAG)
              .first()
              .trigger('entry').length > 0
            ) return false;
          }
          return false;
        }
      });
  });
}

function oalecd8e_switchChineseSetup() {
  oalecd8e.query('.oalecd8e_chn').hide();

  if (oalecd8e.query('.entry .oalecd8e_chn').length != 0) {
    oalecd8e.query('.oalecd8e_show_all').off('.oalecd8e_lang')
      .on('click.oalecd8e_lang', oalecd8e_switchChineseAll);

    oalecd8e.query(OxfordTagSwitchCNTAG)
      .off('.oalecd8e_lang')
      .on('click.oalecd8e_lang entry.oalecd8e_lang', oalecd8e_switchChinese)
      .css('cursor', 'pointer');
  }
}

var oalecd8e_lastSwitchElement;

function oalecd8e_switchChinese(event) {
  // console.log('switch cn');
  if (this === oalecd8e_lastSwitchElement)
    return;

  if (oalecd8e.query(this).is(OxfordTagSwitchCNSIBLINGS)) {
    oalecd8e.query(this).siblings(OxfordTagChineseTexT).oalecd8e_toggle();
    oalecd8e.query(this).children(OxfordTagSwitchCNSIBLINGS).oalecd8e_toggle();
  } else {
    oalecd8e.query(this).children(OxfordTagChineseTexT).oalecd8e_toggle();
  }

  oalecd8e.query(window).scroll();

  setTimeout(function() {
    oalecd8e_lastSwitchElement = null;
  }, 200);
}

function oalecd8e_switchChineseAll(event) {
  oalecd8e.query(OxfordTagChineseTexT).oalecd8e_toggle();
  if (oalecd8e.query('.oalecd8e_show_all.active').length != 0) {
    oalecd8e.query('.oalecd8e_show_all').removeClass("active");
  } else {
    oalecd8e.query('.oalecd8e_show_all').addClass("active");
  }
}

var oalecd8e_slideDuration = 300;

function oalecd8e_extendJQuery() {
  oalecd8e.query.fn.extend({
    oalecd8e_show: function() {
      return this.each(function() {
        if (typeof(oalecd8e.query.fn.fadeIn) == "undefined") {
          oalecd8e.query(this).show();
        } else {
          if (oalecd8e.query(this).css('display') == 'block') {
            oalecd8e.query(this)
              .fadeIn({
                duration: oalecd8e_slideDuration,
                queue: false
              })
              .slideDown(oalecd8e_slideDuration);
          } else {
            oalecd8e.query(this)
              .fadeIn({
                duration: oalecd8e_slideDuration
              });
          }
        }
      });
    },
    oalecd8e_hide: function() {
      return this.each(function() {
        if (typeof(oalecd8e.query.fn.fadeOut) == "undefined") {
          oalecd8e.query(this).hide();
        } else {
          if (oalecd8e.query(this).css('display') == 'block') {
            oalecd8e.query(this)
              .fadeOut({
                duration: oalecd8e_slideDuration,
                queue: false
              })
              .slideUp(oalecd8e_slideDuration);
          } else {
            oalecd8e.query(this)
              .fadeOut({
                duration: oalecd8e_slideDuration,
                queue: false
              });
          }
        }
      });
    },
    oalecd8e_toggle: function(option) {
      return this.each(function(index, element) {
        if ((typeof(option) != 'undefined') ? option : !oalecd8e.query(this).is(":visible")) {
          oalecd8e.query(this).oalecd8e_show();
        } else {
          oalecd8e.query(this).oalecd8e_hide();
        }
      });
    },
    oalecd8e_slideToggle: function(option) {
      return this.each(function(index, element) {
        if (typeof(oalecd8e.query.fn.slideToggle) == "undefined") {
          oalecd8e.query(this).toggle(option);
        } else {
          oalecd8e.query(this).slideToggle(option);
        }
      });
    },
    oalecd8e_fadeIn: function(option) {
      return this.each(function(index, element) {
        if (typeof(oalecd8e.query.fn.fadeIn) == "undefined") {
          oalecd8e.query(this).show(option);
        } else {
          oalecd8e.query(this).fadeIn(option);
        }
      });
    },
    oalecd8e_fadeOut: function(option) {
      return this.each(function(index, element) {
        if (typeof(oalecd8e.query.fn.fadeOut) == "undefined") {
          oalecd8e.query(this).hide(option);
        } else {
          oalecd8e.query(this).fadeOut(option);
        }
      });
    }
  });
}
