<h4>Sample Notex</h4>
<p>
Any line that does not start with one of Notex's special characters will be put in an HTML p tag. If there is more than one newline between two lines of non-special-character-prefixed text, they will be separate p tags.
</p>
<p>
\ The same goes for lines that start with a \ character, but it will place the line verbatim. So if you wanted to write HTML for the end document, or wanted to start the line with a special character then it will be placed in the document.
</p>
<p>
- hyphens are used for unordered lists
</p>
<p>
. and periods for ordered lists
</p>
<p>
{
</p>
<p>
Anything between the braces will be in a subsection
</p>
<p>
 {
</p>
<p>
And they can be infinitely nested
</p>
<p>
}
</p>
<p>
}
</p>
<p>
My favorite feature I've implemented are plugins. There are 3 different plugin scopes:
</p>
<ol>
<li> in-line: will substitute in place in a line </li>
<li> single-line: will substitute the whole line </li>
<li> document: Used for more advanced features, will be fed each token and can do whatever it likes </li>
</ol>
<p>
plugins can be called with /@pluginName:arguments or just /@pluginName, or for elements that effect the head section of the html document !key=@pluginName
</p>
<p>
For example with the theming of the document I have a document scope plugin named "dark" so to set the style I have !style=@dark. When it is called, it determines the deepest subsection level and it creates the relevant css, placing it in the style section in the head. 
</p>
<p>
For images I use the single-line level plugin, "img", which is called by /@img:/path/to/image,extra_args=foo 
</p>
<p>
And for generating timestamps I have the in-line plugin /@now 
</p>
