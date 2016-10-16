package releaseinfo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQualityParsing(t *testing.T) {
	for idx, test := range []struct {
		postTitle string
		quality   Quality
		proper    bool
	}{
		// SDTV
		{"S07E23 .avi ", QualitySDTV, false},
		{"The.Shield.S01E13.x264-CtrlSD", QualitySDTV, false},
		{"Nikita S02E01 HDTV XviD 2HD", QualitySDTV, false},
		{"Gossip Girl S05E11 PROPER HDTV XviD 2HD", QualitySDTV, true},
		{"The Jonathan Ross Show S02E08 HDTV x264 FTP", QualitySDTV, false},
		{"White.Van.Man.2011.S02E01.WS.PDTV.x264-TLA", QualitySDTV, false},
		{"White.Van.Man.2011.S02E01.WS.PDTV.x264-REPACK-TLA", QualitySDTV, true},
		{"The Real Housewives of Vancouver S01E04 DSR x264 2HD", QualitySDTV, false},
		{"Vanguard S01E04 Mexicos Death Train DSR x264 MiNDTHEGAP", QualitySDTV, false},
		{"Chuck S11E03 has no periods or extension HDTV", QualitySDTV, false},
		{"Chuck.S04E05.HDTV.XviD-LOL", QualitySDTV, false},
		{"Sonny.With.a.Chance.S02E15.avi", QualitySDTV, false},
		{"Sonny.With.a.Chance.S02E15.xvid", QualitySDTV, false},
		{"Sonny.With.a.Chance.S02E15.divx", QualitySDTV, false},
		{"The.Girls.Next.Door.S03E06.HDTV-WiDE", QualitySDTV, false},
		{"Degrassi.S10E27.WS.DSR.XviD-2HD", QualitySDTV, false},
		{"[HorribleSubs] Yowamushi Pedal - 32 [480p]", QualitySDTV, false},
		{"[CR] Sailor Moon - 004 [480p][48CE2D0F]", QualitySDTV, false},
		{"[Hatsuyuki] Naruto Shippuuden - 363 [848x480][ADE35E38]", QualitySDTV, false},
		{"Muppet.Babies.S03.TVRip.XviD-NOGRP", QualitySDTV, false},

		// DVD
		{"WEEDS.S03E01-06.DUAL.XviD.Bluray.AC3-REPACK.-HELLYWOOD.avi", QualityDVD, true},
		{"The.Shield.S01E13.NTSC.x264-CtrlSD", QualityDVD, false},
		{"WEEDS.S03E01-06.DUAL.BDRip.XviD.AC3.-HELLYWOOD", QualityDVD, false},
		{"WEEDS.S03E01-06.DUAL.BDRip.X-viD.AC3.-HELLYWOOD", QualityDVD, false},
		{"WEEDS.S03E01-06.DUAL.BDRip.AC3.-HELLYWOOD", QualityDVD, false},
		{"WEEDS.S03E01-06.DUAL.BDRip.XviD.AC3.-HELLYWOOD.avi", QualityDVD, false},
		{"WEEDS.S03E01-06.DUAL.XviD.Bluray.AC3.-HELLYWOOD.avi", QualityDVD, false},
		{"The.Girls.Next.Door.S03E06.DVDRip.XviD-WiDE", QualityDVD, false},
		{"The.Girls.Next.Door.S03E06.DVD.Rip.XviD-WiDE", QualityDVD, false},
		{"the.shield.1x13.circles.ws.xvidvd-tns", QualityDVD, false},
		{"the_x-files.9x18.sunshine_days.ac3.ws_dvdrip_xvid-fov.avi", QualityDVD, false},
		{"[FroZen] Miyuki - 23 [DVD][7F6170E6]", QualityDVD, false},
		{"Hannibal.S01E05.576p.BluRay.DD5.1.x264-HiSD", QualityDVD, false},
		{"Hannibal.S01E05.480p.BluRay.DD5.1.x264-HiSD", QualityDVD, false},
		{"Heidi Girl of the Alps (BD)(640x480(RAW) (BATCH 1) (1-13)", QualityDVD, false},
		{"[Doki] Clannad - 02 (848x480 XviD BD MP3) [95360783]", QualityDVD, false},

		// WebDL 480P
		{"Elementary.S01E10.The.Leviathan.480p.WEB-DL.x264-mSD", QualityWEBDL480p, false},
		{"Glee.S04E10.Glee.Actually.480p.WEB-DL.x264-mSD", QualityWEBDL480p, false},
		{"The.Big.Bang.Theory.S06E11.The.Santa.Simulation.480p.WEB-DL.x264-mSD", QualityWEBDL480p, false},
		{"Da.Vincis.Demons.S02E04.480p.WEB.DL.nSD.x264-NhaNc3", QualityWEBDL480p, false},

		// HDTV 720P
		{"Dexter - S01E01 - Title [HDTV]", QualityHDTV720p, false},
		{"Dexter - S01E01 - Title [HDTV-720p]", QualityHDTV720p, false},
		{"Pawn Stars S04E87 REPACK 720p HDTV x264 aAF", QualityHDTV720p, true},
		{"Sonny.With.a.Chance.S02E15.720p", QualityHDTV720p, false},
		{"S07E23 - [HDTV-720p].mkv ", QualityHDTV720p, false},
		{"Chuck - S22E03 - MoneyBART - HD TV.mkv", QualityHDTV720p, false},
		{"S07E23.mkv ", QualityHDTV720p, false},
		{"Two.and.a.Half.Men.S08E05.720p.HDTV.X264-DIMENSION", QualityHDTV720p, false},
		{"Sonny.With.a.Chance.S02E15.mkv", QualityHDTV720p, false},
		{`"E:\Downloads\tv\The.Big.Bang.Theory.S01E01.720p.HDTV\ajifajjjeaeaeqwer_eppj.avi`, QualityHDTV720p, false},
		{"Gem.Hunt.S01E08.Tourmaline.Nepal.720p.HDTV.x264-DHD", QualityHDTV720p, false},
		{"[Underwater-FFF] No Game No Life - 01 (720p) [27AAA0A0]", QualityHDTV720p, false},
		{"[Doki] Mahouka Koukou no Rettousei - 07 (1280x720 Hi10P AAC) [80AF7DDE]", QualityHDTV720p, false},
		{"[Doremi].Yes.Pretty.Cure.5.Go.Go!.31.[1280x720].[C65D4B1F].mkv", QualityHDTV720p, false},
		{"[HorribleSubs]_Fairy_Tail_-_145_[720p]", QualityHDTV720p, false},
		{"[Eveyuu] No Game No Life - 10 [Hi10P 1280x720 H264][10B23BD8]", QualityHDTV720p, false},
		{"Hells.Kitchen.US.S12E17.HR.WS.PDTV.X264-DIMENSION", QualityHDTV720p, false},
		{"Survivorman.The.Lost.Pilots.Summer.HR.WS.PDTV.x264-DHD", QualityHDTV720p, false},

		// HDTV 1080P
		{"Under the Dome S01E10 Let the Games Begin 1080p", QualityHDTV1080p, false},
		{"DEXTER.S07E01.ARE.YOU.1080P.HDTV.X264-QCF", QualityHDTV1080p, false},
		{"DEXTER.S07E01.ARE.YOU.1080P.HDTV.x264-QCF", QualityHDTV1080p, false},
		{"DEXTER.S07E01.ARE.YOU.1080P.HDTV.proper.X264-QCF", QualityHDTV1080p, true},
		{"Dexter - S01E01 - Title [HDTV-1080p]", QualityHDTV1080p, false},
		{"[HorribleSubs] Yowamushi Pedal - 32 [1080p]", QualityHDTV1080p, false},

		// WebDL 720P
		{"Arrested.Development.S04E01.720p.WEBRip.AAC2.0.x264-NFRiP", QualityWEBDL720p, false},
		{"Vanguard S01E04 Mexicos Death Train 720p WEB DL", QualityWEBDL720p, false},
		{"Hawaii Five 0 S02E21 720p WEB DL DD5 1 H 264", QualityWEBDL720p, false},
		{"Castle S04E22 720p WEB DL DD5 1 H 264 NFHD", QualityWEBDL720p, false},
		{"Chuck - S11E06 - D-Yikes! - 720p WEB-DL.mkv", QualityWEBDL720p, false},
		{"Sonny.With.a.Chance.S02E15.720p.WEB-DL.DD5.1.H.264-SURFER", QualityWEBDL720p, false},
		{"S07E23 - [WEBDL].mkv ", QualityWEBDL720p, false},
		{"Fringe S04E22 720p WEB-DL DD5.1 H264-EbP.mkv", QualityWEBDL720p, false},
		{"House.S04.720p.Web-Dl.Dd5.1.h264-P2PACK", QualityWEBDL720p, false},
		{"Da.Vincis.Demons.S02E04.720p.WEB.DL.nSD.x264-NhaNc3", QualityWEBDL720p, false},
		{"CSI.Miami.S04E25.720p.iTunesHD.AVC-TVS", QualityWEBDL720p, false},
		{"Castle.S06E23.720p.WebHD.h264-euHD", QualityWEBDL720p, false},
		{"The.Nightly.Show.2016.03.14.720p.WEB.x264-spamTV", QualityWEBDL720p, false},
		{"The.Nightly.Show.2016.03.14.720p.WEB.h264-spamTV", QualityWEBDL720p, false},

		// WebDL 1080P
		{"Arrested.Development.S04E01.iNTERNAL.1080p.WEBRip.x264-QRUS", QualityWEBDL1080p, false},
		{"CSI NY S09E03 1080p WEB DL DD5 1 H264 NFHD", QualityWEBDL1080p, false},
		{"Two and a Half Men S10E03 1080p WEB DL DD5 1 H 264 NFHD", QualityWEBDL1080p, false},
		{"Criminal.Minds.S08E01.1080p.WEB-DL.DD5.1.H264-NFHD", QualityWEBDL1080p, false},
		{"Its.Always.Sunny.in.Philadelphia.S08E01.1080p.WEB-DL.proper.AAC2.0.H.264", QualityWEBDL1080p, true},
		{"Two and a Half Men S10E03 1080p WEB DL DD5 1 H 264 REPACK NFHD", QualityWEBDL1080p, true},
		{"Glee.S04E09.Swan.Song.1080p.WEB-DL.DD5.1.H.264-ECI", QualityWEBDL1080p, false},
		{"The.Big.Bang.Theory.S06E11.The.Santa.Simulation.1080p.WEB-DL.DD5.1.H.264", QualityWEBDL1080p, false},
		{"Rosemary's.Baby.S01E02.Night.2.[WEBDL-1080p].mkv", QualityWEBDL1080p, false},
		{"The.Nightly.Show.2016.03.14.1080p.WEB.x264-spamTV", QualityWEBDL1080p, false},
		{"The.Nightly.Show.2016.03.14.1080p.WEB.h264-spamTV", QualityWEBDL1080p, false},
		{"Psych.S01.1080p.WEB-DL.AAC2.0.AVC-TrollHD", QualityWEBDL1080p, false},
		{"Series Title S06E08 1080p WEB h264-EXCLUSIVE", QualityWEBDL1080p, false},
		{"Series Title S06E08 No One PROPER 1080p WEB DD5 1 H 264-EXCLUSIVE", QualityWEBDL1080p, true},
		{"Series Title S06E08 No One PROPER 1080p WEB H 264-EXCLUSIVE", QualityWEBDL1080p, true},
		{"The.Simpsons.S25E21.Pay.Pal.1080p.WEB-DL.DD5.1.H.264-NTb", QualityWEBDL1080p, false},

		// WebDL 2160P
		{"CASANOVA S01E01.2160P AMZN WEBRIP DD2.0 HI10P X264-TROLLUHD", QualityWEBDL2160p, false},
		{"JUST ADD MAGIC S01E01.2160P AMZN WEBRIP DD2.0 X264-TROLLUHD", QualityWEBDL2160p, false},
		{"The.Man.In.The.High.Castle.S01E01.2160p.AMZN.WEBRip.DD2.0.Hi10p.X264-TrollUHD", QualityWEBDL2160p, false},
		{"The Man In the High Castle S01E01 2160p AMZN WEBRip DD2.0 Hi10P x264-TrollUHD", QualityWEBDL2160p, false},
		{"The.Nightly.Show.2016.03.14.2160p.WEB.x264-spamTV", QualityWEBDL2160p, false},
		{"The.Nightly.Show.2016.03.14.2160p.WEB.h264-spamTV", QualityWEBDL2160p, false},
		{"The.Nightly.Show.2016.03.14.2160p.WEB.PROPER.h264-spamTV", QualityWEBDL2160p, true},

		// Bluray 720P
		{"WEEDS.S03E01-06.DUAL.Bluray.AC3.-HELLYWOOD.avi", QualityBluray720p, false},
		{"Chuck - S01E03 - Come Fly With Me - 720p BluRay.mkv", QualityBluray720p, false},
		{"The Big Bang Theory.S03E01.The Electric Can Opener Fluctuation.m2ts", QualityBluray720p, false},
		{"Revolution.S01E02.Chained.Heat.[Bluray720p].mkv", QualityBluray720p, false},
		{"[FFF] DATE A LIVE - 01 [BD][720p-AAC][0601BED4]", QualityBluray720p, false},
		{"[coldhell] Pupa v3 [BD720p][03192D4C]", QualityBluray720p, false},
		{"[RandomRemux] Nobunagun - 01 [720p BD][043EA407].mkv", QualityBluray720p, false},
		{"[Kaylith] Isshuukan Friends Specials - 01 [BD 720p AAC][B7EEE164].mkv", QualityBluray720p, false},
		{"WEEDS.S03E01-06.DUAL.Blu-ray.AC3.-HELLYWOOD.avi", QualityBluray720p, false},
		{"WEEDS.S03E01-06.DUAL.720p.Blu-ray.AC3.-HELLYWOOD.avi", QualityBluray720p, false},
		{"[Elysium]Lucky.Star.01(BD.720p.AAC.DA)[0BB96AD8].mkv", QualityBluray720p, false},
		{"Battlestar.Galactica.S01E01.33.720p.HDDVD.x264-SiNNERS.mkv", QualityBluray720p, false},
		{"The.Expanse.S01E07.RERIP.720p.BluRay.x264-DEMAND", QualityBluray720p, true},

		// Bluray 1080p
		{"Chuck - S01E03 - Come Fly With Me - 1080p BluRay.mkv", QualityBluray1080p, false},
		{"Sons.Of.Anarchy.S02E13.1080p.BluRay.x264-AVCDVD", QualityBluray1080p, false},
		{"Revolution.S01E02.Chained.Heat.[Bluray1080p].mkv", QualityBluray1080p, false},
		{"[FFF] Namiuchigiwa no Muromi-san - 10 [BD][1080p-FLAC][0C4091AF]", QualityBluray1080p, false},
		{"[coldhell] Pupa v2 [BD1080p][5A45EABE].mkv", QualityBluray1080p, false},
		{"[Kaylith] Isshuukan Friends Specials - 01 [BD 1080p FLAC][429FD8C7].mkv", QualityBluray1080p, false},
		{"[Zurako] Log Horizon - 01 - The Apocalypse (BD 1080p AAC) [7AE12174].mkv", QualityBluray1080p, false},
		{"WEEDS.S03E01-06.DUAL.1080p.Blu-ray.AC3.-HELLYWOOD.avi", QualityBluray1080p, false},
		{"[Coalgirls]_Durarara!!_01_(1920x1080_Blu-ray_FLAC)_[8370CB8F].mkv", QualityBluray1080p, false},

		// RAWHD
		{"POI S02E11 1080i HDTV DD5.1 MPEG2-TrollHD", QualityRAWHD, false},
		{"How I Met Your Mother S01E18 Nothing Good Happens After 2 A.M. 720p HDTV DD5.1 MPEG2-TrollHD", QualityRAWHD, false},
		{"The Voice S01E11 The Finals 1080i HDTV DD5.1 MPEG2-TrollHD", QualityRAWHD, false},
		{"Californication.S07E11.1080i.HDTV.DD5.1.MPEG2-NTb.ts", QualityRAWHD, false},
		{"Game of Thrones S04E10 1080i HDTV MPEG2 DD5.1-CtrlHD.ts", QualityRAWHD, false},
		{"VICE.S02E05.1080i.HDTV.DD2.0.MPEG2-NTb.ts", QualityRAWHD, false},
		{"Show - S03E01 - Episode Title Raw-HD.ts", QualityRAWHD, false},
		{"Saturday.Night.Live.Vintage.S10E09.Eddie.Murphy.The.Honeydrippers.1080i.UPSCALE.HDTV.DD5.1.MPEG2-zebra", QualityRAWHD, false},
		{"The.Colbert.Report.2011-08-04.1080i.HDTV.MPEG-2-CtrlHD", QualityRAWHD, false},

		// Unknown
		{"Sonny.With.a.Chance.S02E15", QualityUnknown, false},
		{"Law & Order: Special Victims Unit - 11x11 - Quickie", QualityUnknown, false},
		{"Series.Title.S01E01.webm", QualityUnknown, false},
		{"Droned.S01E01.The.Web.MT-dd", QualityUnknown, false},
	} {
		result := ParseQuality(test.postTitle)

		require.Equal(t, test.quality, result.Quality,
			fmt.Sprintf("Row %d should have the correct quality", idx+1))

		if test.proper {
			require.Equal(t, result.Revision, RevisionProper,
				fmt.Sprintf("Row %d should be version 2 (proper)", idx+1))
		} else {
			require.Equal(t, result.Revision, RevisionDefault,
				fmt.Sprintf("Row %d should be version 1", idx+1))
		}
	}
}

func TestParsingOurOwnQualityNames(t *testing.T) {
	for idx, q := range []Quality{
		QualitySDTV,
		QualityDVD,
		QualityWEBDL480p,
		QualityHDTV720p,
		QualityHDTV1080p,
		QualityHDTV2160p,
		QualityWEBDL720p,
		QualityWEBDL1080p,
		QualityWEBDL2160p,
		QualityBluray720p,
		QualityBluray1080p,
		QualityBluray2160p,
	} {
		fileName := fmt.Sprintf("My series S01E01 [%q]", q.Name)
		result := ParseQuality(fileName)
		require.Equal(t, q, result.Quality,
			fmt.Sprintf("Row %d should parse %s out of %s", idx+1, q.Name, fileName))
	}
}

func TestQualitySourceIsSetWhenParsing(t *testing.T) {
	for idx, test := range []struct {
		postTitle string
		source    string
	}{
		{"Saturday.Night.Live.Vintage.S10E09.Eddie.Murphy.The.Honeydrippers.1080i.UPSCALE.HDTV.DD5.1.MPEG2-zebra", "name"},
		{"Dexter - S01E01 - Title [HDTV-1080p]", "name"},
		{"[CR] Sailor Moon - 004 [480p][48CE2D0F]", "name"},
		{"White.Van.Man.2011.S02E01.WS.PDTV.x264-REPACK-TLA", "name"},
		{"Revolution.S01E02.Chained.Heat.mkv", "extension"},
		{"Dexter - S01E01 - Title.avi", "extension"},
		{"the_x-files.9x18.sunshine_days.avi", "extension"},
		{"[CR] Sailor Moon - 004 [48CE2D0F].avi", "extension"},
	} {
		result := ParseQuality(test.postTitle)

		require.Equal(t, test.source, result.QualitySource,
			fmt.Sprintf("Row %d should have the correct quality source", idx+1))
	}
}

func TestParsingRevision(t *testing.T) {
	for idx, test := range []struct {
		postTitle string
		revision  int
	}{
		{"Chuck.S04E05.HDTV.XviD-LOL", 0},
		{"Gold.Rush.S04E05.Garnets.or.Gold.REAL.REAL.PROPER.HDTV.x264-W4F", 3},
		{"Chuck.S03E17.REAL.PROPER.720p.HDTV.x264-ORENJI-RP", 2},
		{"Covert.Affairs.S05E09.REAL.PROPER.HDTV.x264-KILLERS", 2},
		{"Mythbusters.S14E01.REAL.PROPER.720p.HDTV.x264-KILLERS", 2},
		{"Orange.Is.the.New.Black.s02e06.real.proper.720p.webrip.x264-2hd", 2},
		{"Top.Gear.S21E07.Super.Duper.Real.Proper.HDTV.x264-FTP", 2},
		{"Top.Gear.S21E07.PROPER.HDTV.x264-RiVER-RP", 1},
		{"House.S07E11.PROPER.REAL.RERIP.1080p.BluRay.x264-TENEIGHTY", 2},
		{"[MGS] - Kuragehime - Episode 02v2 - [D8B6C90D]", 1},
		{"[Hatsuyuki] Tokyo Ghoul - 07 [v2][848x480][23D8F455].avi", 1},
		{"[DeadFish] Barakamon - 01v3 [720p][AAC]", 2},
		{"[DeadFish] Momo Kyun Sword - 01v4 [720p][AAC]", 3},
		{"The Real Housewives of Some Place - S01E01 - Why are we doing this?", 0},
		{"[Vivid-Asenshi] Akame ga Kill - 04v2 [266EE983]", 1},
		{"[Vivid-Asenshi] Akame ga Kill - 03v2 [66A05817]", 1},
		{"[Vivid-Asenshi] Akame ga Kill - 02v2 [1F67AB55]", 1},
	} {
		result := ParseQuality(test.postTitle)

		require.Equal(t, test.revision, result.Revision,
			fmt.Sprintf("Row %d should have the correct revision", idx+1))
	}
}
