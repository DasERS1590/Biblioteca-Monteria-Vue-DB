export function getIdUser(){
  const usder = getUser();
  return usder?.id;
}

export function getUser() {
  const userStr = localStorage.getItem("user");
  if (!userStr || userStr === "undefined") return null;
  try {
    return JSON.parse(userStr);
  } catch {
    return null;
  }
}

export function getToken() {
  const user = getUser();
  return user?.token;
} 

