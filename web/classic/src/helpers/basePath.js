function normalizeBasePath(value) {
  if (typeof value !== 'string') return '';
  const trimmed = value.trim();
  if (!trimmed || trimmed === '/') return '';
  const withLeadingSlash = trimmed.startsWith('/') ? trimmed : `/${trimmed}`;
  return withLeadingSlash.replace(/\/+$/, '');
}

function isExternalUrl(value) {
  return /^(?:[a-z][a-z\d+\-.]*:)?\/\//i.test(value);
}

export const APP_BASE_PATH = normalizeBasePath(
  typeof window !== 'undefined' ? window.__APP_BASE_PATH__ : '',
);

export function withBasePath(path) {
  if (!APP_BASE_PATH || !path || isExternalUrl(path) || path.startsWith('#')) {
    return path;
  }
  if (path === APP_BASE_PATH || path.startsWith(`${APP_BASE_PATH}/`)) {
    return path;
  }
  if (path.startsWith('/')) {
    return `${APP_BASE_PATH}${path}`;
  }
  return path;
}

export function originWithBasePath() {
  const origin = typeof window !== 'undefined' ? window.location.origin : '';
  return `${origin}${APP_BASE_PATH}`;
}

export function absoluteUrlWithBasePath(path) {
  if (isExternalUrl(path)) return path;
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${originWithBasePath()}${normalizedPath}`;
}
