// This file contains utility api calling functions for working with files.
// All the functions adopt go's error handling mechanism, they return
// in this format : const [expectedValue, error] = aFunctionCall();

interface FileInformation {
  title: string;
  path: string;
}
type Result<T, E = Error> = [T, null] | [null, E];

export async function getAllTemplateFiles(): Promise<
  Result<FileInformation[]>
> {
  try {
    const response = await fetch("/files");
    if (!response.ok) {
      throw new Error(`Failed to fetch files: ${response.status}`);
    }
    const data: FileInformation[] = await response.json();
    return [data, null];
  } catch (error) {
    if (error instanceof Error) {
      return [null, error];
    }
    return [null, new Error("Unknown error")];
  }
}

export async function getFileContent(
  filename: string,
): Promise<Result<string>> {
  try {
    const response = await fetch(`/read?file=${filename}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch file content: ${response.status}`);
    }
    const data = await response.text();
    return [data, null];
  } catch (error) {
    if (error instanceof Error) {
      return [null, error];
    }
    return [null, new Error("Unknown error")];
  }
}

export async function saveFileContent(
  filename: string,
  content: string,
): Promise<Result<boolean>> {
  try {
    const response = await fetch(`/write?file=${filename}`, {
      method: "POST",
      headers: {
        "Content-Type": "text/plain",
      },
      body: content,
    });
    if (!response.ok) {
      return [
        null,
        new Error(`Failed to save file content: ${response.status}`),
      ];
    }
    return [response.ok, null];
  } catch (error) {
    if (error instanceof Error) {
      return [null, error];
    }
    return [null, new Error("Unknown error")];
  }
}
