import axios from 'axios';

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken');

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

let isRefreshing = false;
let failedQueue = [];

function processQueue(error, token = null) {
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error);
    } else {
      promise.resolve(token);
    }
  });

  failedQueue = [];
}

api.interceptors.response.use(
    (response) => response,

    async (error) => {
      const originalRequest = error.config;

      if (
          error.response?.status === 401 &&
          !originalRequest._retry
      ) {
        if (isRefreshing) {
          return new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject });
          }).then((token) => {
            originalRequest.headers.Authorization =
                `Bearer ${token}`;

            return api(originalRequest);
          });
        }

        originalRequest._retry = true;
        isRefreshing = true;

        try {
          const response = await axios.post(
              `${import.meta.env.VITE_API_BASE_URL || '/api'}/auth/refresh`,
              {},
              {
                withCredentials: true,
              }
          );

          const newAccessToken =
              response.data.access_token;

          localStorage.setItem(
              'accessToken',
              newAccessToken
          );

          processQueue(null, newAccessToken);

          originalRequest.headers.Authorization =
              `Bearer ${newAccessToken}`;

          return api(originalRequest);
        } catch (refreshError) {
          processQueue(refreshError, null);

          localStorage.removeItem('accessToken');

          window.location.href = '/';

          return Promise.reject(refreshError);
        } finally {
          isRefreshing = false;
        }
      }

      return Promise.reject(error);
    }
);

export async function login(credentials) {
  const response = await api.post('/auth/login', credentials);
  return response.data;
}

export async function register(payload) {
  const response = await api.post('/auth/register', payload);
  return response.data;
}

export async function fillProfile(payload) {
  const response = await api.post('/users/me/profile', payload);
  return response.data;
}

export async function mockUploadAvatar(file) {
  await new Promise((resolve) => {
    window.setTimeout(resolve, 350);
  });

  const safeName = file.name.toLowerCase().replace(/[^a-z0-9.]+/g, '-');
  return `/mock-avatars/${Date.now()}-${safeName}`;
}

export async function getMe() {
  const response = await api.get('/users/me');
  return response.data;
}

export async function getUserProfile(userId) {
  const response = await api.get(`/users/${userId}`);
  return response.data;
}

export async function getUserPosts(userId) {
  const response = await api.get(`/users/${userId}/posts`);
  return response.data;
}

export async function getMyPublishedPosts() {
  const response = await api.get('/users/me/posts/published');
  return response.data;
}

export async function getFeed() {
  const response = await api.get('/feed');
  return response.data;
}

export async function createPost(payload) {
  const response = await api.post('/posts', payload);
  return response.data;
}

export async function verifyCode({ code, token }) {
  const response = await api.post('/auth/verify-email', {
    code,
    token,
  });

  return response.data;
}

export async function resendCode({ email }) {
  const response = await api.post('/auth/resend-verification', {
    email,
  });

  return response.data;
}
export async function createLike(postId) {
  const response = await api.post(`/posts/${postId}/like`);
  return response.data;
}

export async function deleteLike(postId) {
  const response = await api.delete(`/posts/${postId}/like`);
  return response.data;
}

export async function createRepost(postId) {
  const response = await api.post(`/posts/${postId}/repost`);
  return response.data;
}

export async function deleteRepost(postId) {
  const response = await api.delete(`/posts/${postId}/repost`);
  return response.data;
}
export async function getMyDraftPosts() {
  const response = await api.get(
      '/users/me/posts/draft'
  );

  return response.data;
}

export async function searchByName(name) {
  const response = await api.get(
      `/users/search/by-name?name=${encodeURIComponent(name)}`
  );

  return response.data;
}

export async function searchByUsername(username) {
  const response = await api.get(
      `/users/search?username=${encodeURIComponent(username)}`
  );

  return response.data;
}

export async function logout() {
  const response = await api.post(
      '/auth/refresh/logout'
  );

  return response.data;
}
export async function followUser(userId) {
  const response = await api.post(
      `/users/${userId}/follow`
  );

  return response.data;
}

export async function unfollowUser(userId) {
  const response = await api.delete(
      `/users/${userId}/follow`
  );

  return response.data;
}
export async function getPostComments(postId) {
  const response = await api.get(
      `/posts/${postId}/comments`
  );

  return response.data;
}

export async function createComment(
    postId,
    payload
) {
  const response = await api.post(
      `/posts/${postId}/comments`,
      payload
  );

  return response.data;
}

export async function createCommentLike(commentId) {
  const response = await api.post(
      `/comments/${commentId}/like`
  );

  return response.data;
}

export async function deleteCommentLike(commentId) {
  const response = await api.delete(
      `/comments/${commentId}/like`
  );

  return response.data;
}

export async function getNotifications() {
  const response = await api.get('/notifications');

  return response.data;
}
export async function getConversations() {
  const response = await api.get(
      '/conversations'
  );

  return response.data;
}