
import { expect } from 'chai'

import YDDict from '../../../src/util/yd'

const yd = new YDDict()
const def = yd.lookup('hello')
expect(def).to.be.a('string')

describe('YDDict', function () {
  // beforeEach(utils.beforeEach)
  // afterEach(utils.afterEach)
  it('search dict', function () {
    const yd = new YDDict()
    const def = yd.lookup('hello')
    console.log(def)
    expect(def).to.be.a('string')
  })
})
