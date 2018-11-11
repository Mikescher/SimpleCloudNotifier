package com.blackforestbytes.simplecloudnotifier.lib.string;

// from MonoSAMFramework.Portable.DebugTools.CompactJsonFormatter
public class CompactJsonFormatter
{
	private static final String INDENT_STRING = "  ";
	
	public static String formatJSON(String str, int maxIndent)
	{
		int indent = 0;
		boolean quoted = false;
		StringBuilder sb = new StringBuilder();
		char last = ' ';
		for (int i = 0; i < str.length(); i++)
		{
			char ch = str.charAt(i);
			switch (ch)
			{
				case '\r':
				case '\n':
					break;
				case '{':
				case '[':
					sb.append(ch);
					last = ch;
					if (!quoted)
					{
						indent++;
						if (indent >= maxIndent) break;
						sb.append("\n");
						for (int ix = 0; ix < indent; ix++) sb.append(INDENT_STRING);
						last = ' ';
					}
					break;
				case '}':
				case ']':
					if (!quoted)
					{
						indent--;
						if (indent + 1 >= maxIndent) { sb.append(ch); break; }
						sb.append("\n");
						for (int ix = 0; ix < indent; ix++) sb.append(INDENT_STRING);
					}
					sb.append(ch);
					last = ch;
					break;
				case '"':
					sb.append(ch);
					last = ch;
					boolean escaped = false;
					int index = i;
					while (index > 0 && str.charAt(--index) == '\\')
						escaped = !escaped;
					if (!escaped)
						quoted = !quoted;
					break;
				case ',':
					sb.append(ch);
					last = ch;
					if (!quoted)
					{
						if (indent >= maxIndent) { sb.append(' '); last = ' '; break; }
						sb.append("\n");
						for (int ix = 0; ix < indent; ix++) sb.append(INDENT_STRING);
					}
					break;
				case ':':
					sb.append(ch);
					last = ch;
					if (!quoted) { sb.append(" "); last = ' '; }
					break;
				case ' ':
				case '\t':
					if (quoted)
					{
						sb.append(ch);
						last = ch;
					}
					else if (last != ' ')
					{
						sb.append(' ');
						last = ' ';
					}
					break;
				default:
					sb.append(ch);
					last = ch;
					break;
			}
		}
		return sb.toString();
	}

	public static String compressJson(String str, int compressionLevel)
	{
		int indent = 0;
		boolean quoted = false;
		StringBuilder sb = new StringBuilder();
		char last = ' ';
		int compress = 0;
		for (int i = 0; i < str.length(); i++)
		{
			char ch = str.charAt(i);
			switch (ch)
			{
				case '\r':
				case '\n':
					break;
				case '{':
				case '[':

					sb.append(ch);
					last = ch;
					if (!quoted)
					{
						if (compress == 0 && getJsonDepth(str, i) <= compressionLevel)
							compress = 1;
						else if (compress > 0)
							compress++;

						indent++;
						if (compress > 0) break;
						sb.append("\n");
						for (int ix = 0; ix < indent; ix++) sb.append(INDENT_STRING);
						last = ' ';
					}
					break;
				case '}':
				case ']':
					if (!quoted)
					{
						indent--;
						if (compress > 0) { compress--; sb.append(ch); break; }
						compress--;
						sb.append("\n");
						for (int ix = 0; ix < indent; ix++) sb.append(INDENT_STRING);
					}
					sb.append(ch);
					last = ch;
					break;
				case '"':
					sb.append(ch);
					last = ch;
					boolean escaped = false;
					int index = i;
					while (index > 0 && str.charAt(--index) == '\\')
						escaped = !escaped;
					if (!escaped)
						quoted = !quoted;
					break;
				case ',':
					sb.append(ch);
					last = ch;
					if (!quoted)
					{
						if (compress > 0) { sb.append(' '); last = ' '; break; }
						sb.append("\n");
						for (int ix = 0; ix < indent; ix++) sb.append(INDENT_STRING);
					}
					break;
				case ':':
					sb.append(ch);
					last = ch;
					if (!quoted) { sb.append(" "); last = ' '; }
					break;
				case ' ':
				case '\t':
					if (quoted)
					{
						sb.append(ch);
						last = ch;
					}
					else if (last != ' ')
					{
						sb.append(' ');
						last = ' ';
					}
					break;
				default:
					sb.append(ch);
					last = ch;
					break;
			}
		}
		return sb.toString();
	}

	public static int getJsonDepth(String str, int i)
	{
		int maxindent = 0;
		int indent = 0;
		boolean quoted = false;
		for (; i < str.length(); i++)
		{
			char ch = str.charAt(i);
			switch (ch)
			{
				case '{':
				case '[':
					if (!quoted)
					{
						indent++;
						maxindent = Math.max(indent, maxindent);
					}
					break;

				case '}':
				case ']':
					if (!quoted)
					{
						indent--;
						if (indent <= 0) return maxindent;
					}
					break;

				case '"':
					boolean escaped = false;
					int index = i;
					while (index > 0 && str.charAt(--index) == '\\')
						escaped = !escaped;
					if (!escaped)
						quoted = !quoted;
					break;

				default:
					break;
			}
		}
		return maxindent;
	}
}
