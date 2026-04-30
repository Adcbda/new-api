// ============================================================================
// Affiliate Functions
// ============================================================================

import { absoluteUrlWithBasePath } from '@/lib/base-path'

/**
 * Generate affiliate registration link
 */
export function generateAffiliateLink(affCode: string): string {
  if (typeof window === 'undefined') return ''
  return absoluteUrlWithBasePath(`/register?aff=${affCode}`)
}
