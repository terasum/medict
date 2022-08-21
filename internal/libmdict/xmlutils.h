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

#ifndef MDICT_XMLUTILS_H_
#define MDICT_XMLUTILS_H_

#include <map>

// parse xml header info
std::map<std::string, std::string> parseXMLHeader(std::string dicxml);

#endif
