/*
 * =====================================================================================
 *
 *       Filename:  xmlutils.h
 *
 *    Description:  xml parser utils tool functions
 *
 *        Version:  1.0
 *        Created:  01/25/2019 13:45:05
 *       Revision:  none
 *       Compiler:  gcc
 *
 *         Author:  terasum (terasum@163.com)
 *
 * =====================================================================================
 */
#include <iostream>
#include <map>

std::map<std::string, std::string> parseXMLHeader(std::string dicxml) {
  /// std::map<char*, char*> headTag;

  // scaner loop
  //  bool start = false, end = false;
  bool endCorrect = false;
  char endPrevChar = 0, endChar = 0;

  std::map<std::string, std::string> tagm;
  // end correct with "/>" ?
  if (dicxml.length() > 2) {
    endPrevChar = dicxml.at(dicxml.length() - 2);
    endChar = dicxml.at(dicxml.length() - 1);
    if (endPrevChar == '/' && endChar == '>') {
      endCorrect = true;
    }
  } else {
    return tagm;
  }

  unsigned i = 0;

  // get start index
  for (; i < dicxml.length(); ++i) {
    char c = dicxml.at(i);
    if (c != ' ')
      continue;
    else
      break;
  }

  int r = i;

  // loop to length -2 (for '/>')
  for (; r < dicxml.length() - 2; r++) {
  stringloop:
    std::string tmpK;
    std::string tmpV;
    // split flag
    bool sflag = false;
    bool openQuo = false, closeQuo = false;
    bool finished = false;
    int tj = 0;
    for (int j = 0; r + j <= dicxml.length() - 2; ++j) {
      tj = j;
      char c = dicxml.at(r + j);
      // cout << c;

      bool getKey = !sflag;
      bool getValue = openQuo && !closeQuo;

      // ensure get value station
      // sflag == true means key has already getted
      // and openQuo and !closeQuo means current getting value
      if (!getKey && !getValue && c == ' ' && !(r + j == dicxml.length() - 2)) {
        //				std::cout << tmpK
        //<<":"<<tmpV<<std::endl;
        tagm[tmpK] = tmpV;
        r += j;
        // cout<<"bbb"<<r<<endl;
        goto stringloop;
      } else if (!getKey && !getValue && (r + j == dicxml.length() - 2)) {
        //				std::cout << tmpK
        //<<":"<<tmpV<<std::endl;
        tagm[tmpK] = tmpV;
        r = dicxml.length();  // this is for exit tagloop
        // cout<<"aaa"<<r+j<<endl;
        break;
      }
      // get key station, c cannot be space
      else if (getKey && c == ' ') {
        continue;
        // split station, should not in value getting station
      }
      // store key
      else if (getKey && c != ' ' && c != '=') {
        tmpK += c;
      }
      // end getKey start get value
      else if (getKey && c == '=') {
        sflag = true;
        continue;
        // can be key and value getting station
      }
      // ----- get key end, get value start -----
      else if (!getKey && !getValue && c == '"') {
        openQuo = true;
        continue;
      } else if (!getKey && getValue && c != '"') {
        tmpV += c;
        continue;
      } else if (!getKey && getValue && c == '"') {
        closeQuo = true;
        finished = true;
        continue;
      }
    }
    //		printf(" |(r:%d, j:%d) l: %d, gk: %d, gv %d \n", r, tj,
    // dicxml.length(),!sflag,openQuo && !closeQuo);
  }
  //	printf("585: %c", dicxml.at(585));
  return tagm;
}
