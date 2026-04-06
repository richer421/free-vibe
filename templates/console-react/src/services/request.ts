export interface RequestOptionsDTO extends RequestInit {
  baseUrl?: string
}

export async function request<TResponse>(
  input: string,
  options: RequestOptionsDTO = {},
): Promise<TResponse> {
  const { baseUrl = '', headers, ...restOptions } = options
  const response = await fetch(`${baseUrl}${input}`, {
    ...restOptions,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
  })

  if (!response.ok) {
    throw new Error(`Request failed with status ${response.status}`)
  }

  return (await response.json()) as TResponse
}
