import { describe, expect, it } from 'vitest'
import { defaultLayout, moveModule } from './layout'

describe('moveModule', () => {
  it('moves a module without changing the remaining order', () => {
    expect(moveModule(defaultLayout, 'resources', 1).order).toEqual([
      'overview', 'posts', 'resources', 'ai-lab',
    ])
  })
})
