import ky from "ky";

interface ListPostInput {
  skip?: number;
  limit?: number;
}

async function PostList(data: ListPostInput): Promise<Post[]> {
  const queryParams = new URLSearchParams(data as unknown as Record<string, any>).toString();
  return await ky.get(`http://localhost:9090/post/list/?${queryParams}`).json<Post[]>();
}

interface CreatePostInput {
  title: string;
  content: string;
}

interface Post {
  content: string;
  id: number;
  title: string;
}

async function PostCreate(data: CreatePostInput): Promise<Post> {
  return await ky.post("http://localhost:9090/post/create/", {
    json: data
  }).json<Post>();
}

interface GetPostInput {
  id: number;
  author_id: string;
}

async function PostGet(data: GetPostInput): Promise<Post> {
  const queryParams = new URLSearchParams(data as unknown as Record<string, any>).toString();
  return await ky.get(`http://localhost:9090/post/get/?${queryParams}`).json<Post>();
}

