export function validatePath(path: string): { valid: boolean; error?: string } {
  if (!path) {
    return { valid: false, error: 'Path cannot be empty' }
  }

  if (!path.startsWith('/')) {
    return { valid: false, error: 'Path must be absolute' }
  }

  const protectedPaths = ['/', '/System', '/bin', '/sbin', '/usr', '/etc', '/var']
  if (protectedPaths.includes(path)) {
    return { valid: false, error: 'Cannot delete protected system path' }
  }

  return { valid: true }
}

export function validateNotEmpty(value: string, fieldName: string) {
  if (!value || value.trim() === '') {
    return { valid: false, error: `${fieldName} cannot be empty` }
  }
  return { valid: true }
}
