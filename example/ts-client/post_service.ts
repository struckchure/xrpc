interface ListPostInput {
  skip?: number;
  limit?: number;
}

async function PostList(data: ListPostInput): Promise<Post[]> {
  const queryParams = new URLSearchParams(data as unknown as Record<string, any>).toString();
  const response = await fetch(`http://localhost:9090/post/list/?${queryParams}`);
  return response.json();
}

interface CreatePostInput {
  title: string;
  content: string;
}

interface Post {
  id: number;
  title: string;
  content: string;
}

async function PostCreate(data: CreatePostInput): Promise<Post> {
  const response = await fetch("http://localhost:9090/post/create/", {
    method: "POST",
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
  return response.json();
}

interface GetPostInput {
  id?: number;
  author_id?: string;
}

async function PostGet(data: GetPostInput): Promise<Post> {
  const queryParams = new URLSearchParams(data as unknown as Record<string, any>).toString();
  const response = await fetch(`http://localhost:9090/post/get/?${queryParams}`);
  return response.json();
}

