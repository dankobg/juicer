import { kratos } from './client';

export async function getSession() {
  try {
    const resp = await kratos.toSession();
    return resp.data;
  } catch (error) {
    // 401, 403 -> not logged in
    return null;
  }
}
