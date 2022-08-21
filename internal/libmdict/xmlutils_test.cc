#include <xmlutils.h>

#include <iostream>

using namespace std;

int main() {
  string dicxml =
      "<Dictionary GeneratedByEngineVersion=\"2.0\" "
      "RequiredEngineVersion=\"2.0\" Format=\"Html\" KeyCaseSensitive=\"No\" "
      "StripKey=\"Yes\" Encrypted=\"2\" RegisterBy=\"EMail\" "
      "Description=\"Oxford Advanced Learner’s English-Chinese Dictionary "
      "Eighth edition Based on Langheping&apos;s version Modified by "
      "EarthWorm&lt;br/&gt;Headwords: 41969 &lt;br/&gt; Entries: 109473 "
      "&lt;br/&gt; Version: 3.0.0 &lt;br/&gt;Date: 2018.02.18 &lt;br/&gt; Last "
      "Modified By roamlog&lt;br/&gt;\" Title=\"\" Encoding=\"UTF-8\" "
      "CreationDate=\"2018-2-18\" Compact=\"Yes\" Compat=\"Yes\" "
      "Left2Right=\"Yes\" DataSourceFormat=\"106\" StyleSheet=\"a\"/>";
  map<string, string> headinfo = parseXMLHeader(dicxml);
  //	cout<<headinfo["Encoding"]<<endl;
  assert(headinfo["Encoding"] == "UTF-8");
}
