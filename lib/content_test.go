package lib

import (
	"strings"
	"testing"
)

const testContent = `
<div id="container" role="main"><div  type="chapter" uri="/scriptures/nt/mark/1" schema-version="1" doc-id="scriptures_03990_000_mark_001" doc-version="1" hash="xBLIpA">
<div id="/scriptures/nt/mark/1" uri="/scriptures/nt/mark/1" hash="xBLIpA" class="chapter">
<div id="head" uri="/scriptures/nt/mark/1.head" class="heading">
<h1 pid="8zs4R2CyGU-3otWsmfUnhA" hash="mjRioQ">The Gospel According to <br/><span class="dominant"><a href="f_note">St Mark</a></span></h1><p pid="ijnWfv_220Koq-6wR1uwIQ" hash="A_bhkw" class="titleNumber">Chapter 1</p><p pid="sdmjnXqj20an7M0pddkcgw" hash="u6cp3Q" class="studySummary">Jesus is baptized by John—He preaches the gospel, calls disciples, casts out devils, heals the sick, and cleanses a leper.</p></div>
<div class="bodyBlock"> <p id="1" uri="/scriptures/nt/mark/1.1" pid="kq-BAmvcsUWpiuiFEKd3xQ" hash="-Nx6tg" class="verse">1 The beginning of the <sup>a</sup><a href="f_1a">gospel</a> of Jesus Christ, the Son of God;</p> <p id="2" uri="/scriptures/nt/mark/1.2" pid="AXVsTdPL9ECykupkBNytXg" hash="gt2BGA" class="verse">2 As it is written in the prophets, Behold, I send my <sup>a</sup><a href="f_2a">messenger</a> before thy face, which shall prepare thy way before thee.</p> <p id="3" uri="/scriptures/nt/mark/1.3" pid="WCZbiNCR1Uil44fdCsxnkg" hash="j2XA9g" class="verse">3 The <sup>a</sup><a href="f_3a">voice</a> of one crying in the wilderness, Prepare ye the way of the Lord, make his paths straight.</p> <p id="4" uri="/scriptures/nt/mark/1.4" pid="7Pqsm64qKUGvxYuyywDK7w" hash="pt4j5Q" class="verse">4 John did baptize in the wilderness, and <sup>a</sup><a href="f_4a">preach</a> the <sup>b</sup><a href="f_4b">baptism</a> of <span class="pageBreak" page-number="1242"></span><sup>c</sup><a href="f_4c">repentance</a> for the <sup>d</sup><a href="f_4d">remission</a> of sins.</p> <p id="5" uri="/scriptures/nt/mark/1.5" pid="j7N5RCKzR02kWNpA-UqoRw" hash="1PPg1A" class="verse">5 And there went out unto him all the land of Judæa, and they of Jerusalem, and were all <sup>a</sup><a href="f_5a">baptized</a> of him in the river of Jordan, <sup>b</sup><a href="f_5b">confessing</a> their sins.</p> <p id="6" uri="/scriptures/nt/mark/1.6" pid="ZVWKUdylckyojP_3pDFteQ" hash="P5CQog" class="verse">6 And John was <sup>a</sup><a href="f_6a">clothed</a> with <sup>b</sup><a href="f_6b">camel’s</a> hair, and with a girdle of a skin about his loins; and he did eat <sup>c</sup><a href="f_6c">locusts</a> and wild honey;</p> <p id="7" uri="/scriptures/nt/mark/1.7" pid="Ne6tl_CbT0eZR1UPiJPjBg" hash="PgT2rA" class="verse">7 And preached, saying, There cometh one mightier than I after me, the latchet of whose shoes I am not worthy to stoop down and unloose.</p> <p id="8" uri="/scriptures/nt/mark/1.8" pid="VeEdoDPc306PJBYr5-etTw" hash="fdMOlg" class="verse">8 I indeed have baptized you with water: <sup>a</sup><a href="f_8a">but</a> he shall baptize you with the Holy Ghost.</p> <p id="9" uri="/scriptures/nt/mark/1.9" pid="ZdrnzInVHEa5XWqFbnCR2g" hash="-8GHLg" class="verse">9 And it came to pass in those days, that Jesus came from Nazareth of Galilee, and was <sup>a</sup><a href="f_9a">baptized</a> of John in Jordan.</p> <p id="10" uri="/scriptures/nt/mark/1.10" pid="pizLBKZGZkuCP2fgQmx6Tg" hash="bR3HKw" class="verse">10 And straightway coming <sup>a</sup><a href="f_10a">up</a> out of the water, he saw the heavens opened, and the Spirit like a <sup>b</sup><a href="f_10b">dove</a> descending upon him:</p> <p id="11" uri="/scriptures/nt/mark/1.11" pid="xaaSn2SDpEqfRQa5PBQobw" hash="pFnrlg" class="verse">11 And there came a voice from heaven, <span class="clarityWord">saying,</span> Thou art my beloved Son, in whom I am well pleased.</p> <p id="12" uri="/scriptures/nt/mark/1.12" pid="Drfd2BMX7EW4dN5wqqTmLQ" hash="qwaLEw" class="verse">12 <sup>a</sup><a href="f_12a">And</a> immediately the Spirit driveth him into the <sup>b</sup><a href="f_12b">wilderness</a>.</p> <p id="13" uri="/scriptures/nt/mark/1.13" pid="-NyOIsI4u0CI6f5IImtHaQ" hash="CV6Hwg" class="verse">13 And he was there in the wilderness forty days, tempted of Satan; and was with the wild beasts; and the angels ministered unto him.</p> <p id="14" uri="/scriptures/nt/mark/1.14" pid="2Wp5IF2_t0m__h4h7W5f5A" hash="gvJ2CQ" class="verse">14 Now after that John was put in prison, Jesus came into Galilee, <sup>a</sup><a href="f_14a">preaching</a> the gospel of the kingdom of God,</p> <p id="15" uri="/scriptures/nt/mark/1.15" pid="neY2AAmUB0u7FtbgEmU9ZA" hash="0fF_Aw" class="verse">15 And saying, The <sup>a</sup><a href="f_15a">time</a> is fulfilled, and the <sup>b</sup><a href="f_15b">kingdom</a> of God <sup>c</sup><a href="f_15c">is at hand</a>: <sup>d</sup><a href="f_15d">repent</a> ye, and <sup>e</sup><a href="f_15e">believe</a> the gospel.</p> <p id="16" uri="/scriptures/nt/mark/1.16" pid="kEll_NlY5UyrNJsZwpCRXw" hash="hiZqYQ" class="verse">16 Now as he walked by the sea of Galilee, he saw Simon and Andrew his brother casting a net into the sea: for they were fishers.</p> <p id="17" uri="/scriptures/nt/mark/1.17" pid="gY_Ql3D5Kk-KbUBiJkMJrw" hash="9tLOmA" class="verse">17 And Jesus said unto them, Come ye after me, and I will make you to become <sup>a</sup><a href="f_17a">fishers of men</a>.</p> <p id="18" uri="/scriptures/nt/mark/1.18" pid="aWHR21Xc4EyjddK5BofjXQ" hash="8NZMaA" class="verse">18 And straightway they forsook their nets, and followed him.</p> <p id="19" uri="/scriptures/nt/mark/1.19" pid="cnmqPI7hE0q3C24c4tqYKw" hash="w9W4sQ" class="verse">19 And when he had gone a little further thence, he saw James the <span class="clarityWord">son</span> of Zebedee, and John his brother, who also were in the ship mending their nets.</p> <p id="20" uri="/scriptures/nt/mark/1.20" pid="bM67FADuo02Lv8CHdx_MRA" hash="MnsA-g" class="verse">20 And straightway he called them: and they left their father Zebedee in the ship with the hired servants, and went after him.</p> <p id="21" uri="/scriptures/nt/mark/1.21" pid="caibuVw7jEqxLYGmro8TEg" hash="9-RZiA" class="verse">21 And they went into Capernaum; and straightway on the sabbath day he entered into the synagogue, and <sup>a</sup><a href="f_21a">taught</a>.</p> <p id="22" uri="/scriptures/nt/mark/1.22" pid="f6b0Nb-lckeS3kJ8LW9jQw" hash="U0E5mQ" class="verse">22 And they were astonished at his doctrine: for he taught them as one that had <sup>a</sup><a href="f_22a">authority</a>, and not as the <sup>b</sup><a href="f_22b">scribes</a>.</p> <p id="23" uri="/scriptures/nt/mark/1.23" pid="hjdu5cuWtUemdiJaFlJlJw" hash="LuUJiw" class="verse">23 And there was in their synagogue a man with an <sup>a</sup><a href="f_23a">unclean spirit</a>; and he cried out,</p> <p id="24" uri="/scriptures/nt/mark/1.24" pid="Ei2sDM9p0UWZP_e1_t5TaA" hash="b4wneA" class="verse">24 Saying, Let <span class="clarityWord">us</span> alone; <sup>a</sup><a href="f_24a">what</a> have we to do with thee, thou Jesus of Nazareth? art thou come to destroy us? I know thee who thou art, the <sup>b</sup><a href="f_24b">Holy One</a> of God.</p> <p id="25" uri="/scriptures/nt/mark/1.25" pid="LBnvrUUw0UKPIBKQgFSD0A" hash="tt1SoA" class="verse">25 And Jesus <sup>a</sup><a href="f_25a">rebuked</a> him, saying, Hold thy peace, and come out of him.</p> <span class="pageBreak" page-number="1243"></span> <p id="26" uri="/scriptures/nt/mark/1.26" pid="dsg52e48x0OAwXaylIybbw" hash="JrLPDQ" class="verse">26 And when the unclean spirit had <sup>a</sup><a href="f_26a">torn</a> him, and cried with a loud voice, he came out of him.</p> <p id="27" uri="/scriptures/nt/mark/1.27" pid="WbZdAvcaf0-ueEH9Htt2qw" hash="gEJ7pg" class="verse">27 And they were all amazed, insomuch that they questioned among themselves, saying, What thing is this? what new doctrine <span class="clarityWord">is</span> this? for with <sup>a</sup><a href="f_27a">authority</a> commandeth he even the unclean spirits, and they do obey him.</p> <p id="28" uri="/scriptures/nt/mark/1.28" pid="8-4-GkTHvkms3K1f1ottug" hash="RuxEkw" class="verse">28 And immediately his fame spread abroad throughout all the region round about Galilee.</p> <p id="29" uri="/scriptures/nt/mark/1.29" pid="AjIG_GJT7UaumH2I0EcZMQ" hash="fZduqA" class="verse">29 And forthwith, when they were come out of the synagogue, they entered into the house of Simon and Andrew, with James and John.</p> <p id="30" uri="/scriptures/nt/mark/1.30" pid="skIEAA6Z6EG-6PRgwpmIrA" hash="9nEpBQ" class="verse">30 But Simon’s wife’s mother lay sick of a fever, and <sup>a</sup><a href="f_30a">anon</a> they tell him of her.</p> <p id="31" uri="/scriptures/nt/mark/1.31" pid="6F0Ypf6t_kOHXQFoLi9ftw" hash="da8kgg" class="verse">31 And he came and took her by the hand, and <sup>a</sup><a href="f_31a">lifted</a> her up; and immediately the fever left her, and she ministered unto them.</p> <p id="32" uri="/scriptures/nt/mark/1.32" pid="afirGFIrd0qUO1Gs9FnX8w" hash="1gkyqQ" class="verse">32 And at even, when the sun did set, they brought unto him all that were diseased, and them that were possessed with devils.</p> <p id="33" uri="/scriptures/nt/mark/1.33" pid="_xbYTCkjn02AVgBxpkDEug" hash="Obo7DA" class="verse">33 And all the city was gathered together at the door.</p> <p id="34" uri="/scriptures/nt/mark/1.34" pid="B088O3VHK0WRrJQrJGhOvw" hash="nAEMCw" class="verse">34 And he <sup>a</sup><a href="f_34a">healed</a> many that were sick of divers diseases, and cast out many <sup>b</sup><a href="f_34b">devils</a>; and <sup>c</sup><a href="f_34c">suffered not</a> the devils to speak, because they knew him.</p> <p id="35" uri="/scriptures/nt/mark/1.35" pid="oG3iWBxASE-DI5e5jVY0mg" hash="jets8Q" class="verse">35 And in the morning, rising up a great while before day, he went out, and departed into a solitary place, and there prayed.</p> <p id="36" uri="/scriptures/nt/mark/1.36" pid="fBwUaiwwvEC64mms3lHL8A" hash="-jQwgw" class="verse">36 And Simon and they that were with him followed after him.</p> <p id="37" uri="/scriptures/nt/mark/1.37" pid="P6mqvR5jA0CQUvJLUsEx2g" hash="Zz2kqw" class="verse">37 And when they had found him, they said unto him, All <span class="clarityWord">men</span> seek for thee.</p> <p id="38" uri="/scriptures/nt/mark/1.38" pid="gnr3RBsmjkWPjjEmHBYi4w" hash="OTZbjA" class="verse">38 And he said unto them, Let us go into the next towns, that I may <sup>a</sup><a href="f_38a">preach</a> there also: for therefore came I forth.</p> <p id="39" uri="/scriptures/nt/mark/1.39" pid="OKFIIc8aJUarq_REcem0xA" hash="culV3g" class="verse">39 And he preached in their synagogues throughout all Galilee, and cast out <sup>a</sup><a href="f_39a">devils</a>.</p> <p id="40" uri="/scriptures/nt/mark/1.40" pid="VbuZE7Z0vUCSCbDPa3p1IA" hash="P9jtFg" class="verse">40 And there came a <sup>a</sup><a href="f_40a">leper</a> to him, beseeching him, and kneeling down to him, and saying unto him, If thou wilt, thou canst make me <sup>b</sup><a href="f_40b">clean</a>.</p> <p id="41" uri="/scriptures/nt/mark/1.41" pid="n9X0ph-lfkOCX_M1Xm65wQ" hash="ufw1dA" class="verse">41 And Jesus, moved with <sup>a</sup><a href="f_41a">compassion</a>, put forth <span class="clarityWord">his</span> hand, and touched him, and saith unto him, I will; be thou clean.</p> <p id="42" uri="/scriptures/nt/mark/1.42" pid="bMAN_nL6YE-fEKIm2DY6-g" hash="BahF-A" class="verse">42 And as soon as he had spoken, immediately the leprosy departed from him, and he was cleansed.</p> <p id="43" uri="/scriptures/nt/mark/1.43" pid="3YYmAV99iEWTOwm8VGvALg" hash="Y9w-1Q" class="verse">43 And he <sup>a</sup><a href="f_43a">straitly charged him</a>, and forthwith sent him away;</p> <p id="44" uri="/scriptures/nt/mark/1.44" pid="hMvRUAGkmEiLrV94q6dA1w" hash="_Z_8zA" class="verse">44 And saith unto him, See thou say nothing to any man: but go thy way, shew thyself to the priest, and offer for thy cleansing those things which Moses commanded, for a testimony unto them.</p> <p id="45" uri="/scriptures/nt/mark/1.45" pid="G3ZEdA1ZHEKlI47kQEvCPA" hash="oY0t0Q" class="verse">45 But he went out, and began to publish <span class="clarityWord">it</span> much, and to <sup>a</sup><a href="f_45a">blaze abroad</a> the matter, insomuch that Jesus could no more openly enter into the city, but was without in desert places: and they came to him from every quarter.</p> </div></div></div></div>
`

const testContentFiltered01 = `
<div id="container" role="main"><div  type="chapter" uri="/scriptures/nt/mark/1" schema-version="1" doc-id="scriptures_03990_000_mark_001" doc-version="1" hash="xBLIpA">
<div id="/scriptures/nt/mark/1" uri="/scriptures/nt/mark/1" hash="xBLIpA" class="chapter">
<div id="head" uri="/scriptures/nt/mark/1.head" class="heading">
<h1 pid="8zs4R2CyGU-3otWsmfUnhA" hash="mjRioQ">The Gospel According to <br/><span class="dominant"><a href="f_note">St Mark</a></span></h1><p pid="ijnWfv_220Koq-6wR1uwIQ" hash="A_bhkw" class="titleNumber">Chapter 1</p><p pid="sdmjnXqj20an7M0pddkcgw" hash="u6cp3Q" class="studySummary">Jesus is baptized by John—He preaches the gospel, calls disciples, casts out devils, heals the sick, and cleanses a leper.</p></div>
<div class="bodyBlock"> <p id="1" uri="/scriptures/nt/mark/1.1" pid="kq-BAmvcsUWpiuiFEKd3xQ" hash="-Nx6tg" class="verse">1 The beginning of the <sup>a</sup><a href="f_1a">gospel</a> of Jesus Christ, the Son of God;</p>                         <span class="pageBreak" page-number="1243"></span>                     </div></div></div></div>
`

const testContentFiltered1 = `<p id="1" uri="/scriptures/nt/mark/1.1" pid="kq-BAmvcsUWpiuiFEKd3xQ" hash="-Nx6tg" class="verse">1 The beginning of the <sup>a</sup><a href="f_1a">gospel</a> of Jesus Christ, the Son of God;</p>`

func testSearchResult(t *testing.T, p, r SearchResult) {
	testReference(t, p.Reference, r.Reference)
	if p.Weight != r.Weight {
		t.Errorf("      Weight doesn't match %v != %v", p.Weight, r.Weight)
	}
}

func testParagraph(t *testing.T, z *ContentParser, style ParagraphStyle, verse int) {
	if !z.NextParagraph() {
		t.Errorf("     Parse ended too quickly for style %v, verse %v", style, verse)
	}
	if z.ParagraphStyle() != style {
		t.Errorf("    Wrong paragraph style %v vs. %v", z.ParagraphStyle(), style)
	}
	if z.ParagraphVerse() != verse {
		t.Errorf("    Incorrectly assigned verse %v vs. %v", z.ParagraphVerse(), verse)
	}
}
func testText(t *testing.T, z *ContentParser, style TextStyle, text string) {
	if !z.NextText() {
		t.Errorf("    Paragraph ended too quickly for %v \"%v\"", style, text)
	}
	zstyle := z.TextStyle()
	if zstyle != style {
		t.Errorf("    Wrong text style %v vs. %v", zstyle, style)
	}
	ztext := z.Text()
	if ztext != text {
		t.Errorf("    Wrong text content \"%v\" vs. \"%v\"", ztext, text)
	}
}
func testTextEnd(t *testing.T, z *ContentParser) {
	if z.NextText() {
		t.Error("    Paragraph extended too far")
	}
}

func TestContentFilter(t *testing.T) {
	c := Content(testContent)
	c01 := string(c.Filter([]int{0, 1}))
	if strings.TrimSpace(c01) != strings.TrimSpace(testContentFiltered01) {
		t.Errorf("Filter test failed c01 between:\n%v\n and \n%v", c01, testContentFiltered01)
	}
	c1 := string(c.Filter([]int{1}))
	if strings.TrimSpace(c1) != strings.TrimSpace(testContentFiltered1) {
		t.Errorf("Filter test failed c1 between:\n%v\n and \n%v", c1, testContentFiltered1)
	}
}

func TestContentParse1(t *testing.T) {
	t.Log("Parsing content")
	c := Content(testContent)
	c = c.Filter([]int{1})
	z := c.Parse()

	testParagraph(t, z, ParagraphStyleNormal, 1)
	testText(t, z, TextStyleNormal, "1 The beginning of the ")
	testText(t, z, TextStyleFootnote, "a")
	testText(t, z, TextStyleLink, "gospel")
	testText(t, z, TextStyleNormal, " of Jesus Christ, the Son of God;")
	testTextEnd(t, z)
	if z.NextParagraph() {
		t.Error("    Parse extended too far")
	}
}

func TestContentParse(t *testing.T) {
	t.Log("Parsing content")
	c := Content(testContent)
	c = c.Filter([]int{0, 1})
	z := c.Parse()

	testParagraph(t, z, ParagraphStyleTitle, 0)
	testText(t, z, TextStyleNormal, "The Gospel According to ")
	testText(t, z, TextStyleLink, "St Mark")
	testTextEnd(t, z)
	testParagraph(t, z, ParagraphStyleChapter, 0)
	testText(t, z, TextStyleNormal, "Chapter 1")
	testTextEnd(t, z)
	testParagraph(t, z, ParagraphStyleSummary, 0)
	testText(t, z, TextStyleNormal, "Jesus is baptized by John—He preaches the gospel, calls disciples, casts out devils, heals the sick, and cleanses a leper.")
	testTextEnd(t, z)
	testParagraph(t, z, ParagraphStyleNormal, 1)
	testText(t, z, TextStyleNormal, "The beginning of the ")
	testText(t, z, TextStyleFootnote, "a")
	testText(t, z, TextStyleLink, "gospel")
	testText(t, z, TextStyleNormal, " of Jesus Christ, the Son of God;")
	testTextEnd(t, z)
	if z.NextParagraph() {
		t.Error("    Parse extended too far")
	}
}

func TestContentSearch(t *testing.T) {
	c := Content(testContent)

	test := func(weight int, verses []int, keywords ...string) {
		r := SearchResult{}
		r.Weight = weight
		r.VersesHighlighted = verses
		t.Logf("Testing strings \"%v\" for match %v", keywords, r)
		p := c.Search(keywords)
		testSearchResult(t, r, p)
	}

	test(0, nil, "blockBody")           // Make sure strings in attributes don't match
	test(5, []int{1, 14, 15}, "gospel") // Check strings in verses and headers
}
