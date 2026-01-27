const url = process.env.VUE_APP_API_URL;

export async function query(path, params = null) {
  const query = params ? `?${new URLSearchParams(params)}` : "";

  try {
    const res = await fetch(`${url}/${path}${query}`);
    return await res.json();
  } catch {
    return null;
  }
}

export async function post(path, params = null, body = null) {
  const query = params ? `?${new URLSearchParams(params)}` : "";

  try {
    const res = await fetch(`${url}/${path}${query}`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: body ? JSON.stringify(body) : null,
    });

    return await res.json();
  } catch {
    return null;
  }
}
