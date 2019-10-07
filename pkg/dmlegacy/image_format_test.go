package dmlegacy

import "testing"

func TestParseResize2fsOutputForMinSize(t *testing.T) {
	cases := []struct {
		name, out string
		expected  int64
	}{
		{
			name: `POSIX, C, C.UTF-8, en_AG, en_AG.utf8, en_AU.utf8, en_BW.utf8, en_CA.utf8, en_DK.utf8, en_GB.utf8, en_HK.utf8, en_IE.utf8, en_IL, en_IL.utf8, en_IN, en_IN.utf8, en_NG, en_NG.utf8, en_NZ.utf8, en_PH.utf8, en_SG.utf8, en_US.utf8, en_ZA.utf8, en_ZM, en_ZM.utf8, en_ZW.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Estimated minimum size of the filesystem: 797480`,
			expected: 797480,
		},
		{
			name: `af_ZA.utf8, am_ET, am_ET.utf8, an_ES.utf8, ar_AE.utf8, ar_BH.utf8, ar_DZ.utf8, ar_EG.utf8, ar_IN, ar_IN.utf8, ar_IQ.utf8, ar_JO.utf8, ar_KW.utf8, ar_LB.utf8, ar_LY.utf8, ar_MA.utf8, ar_OM.utf8, ar_QA.utf8, ar_SA.utf8, ar_SD.utf8, ar_SS, ar_SS.utf8, ar_SY.utf8, ar_TN.utf8, ar_YE.utf8, as_IN, as_IN.utf8, ast_ES.utf8, az_AZ, az_AZ.utf8, be_BY.utf8, be_BY.utf8@latin, be_BY@latin, bg_BG.utf8, bn_BD, bn_BD.utf8, bn_IN, bn_IN.utf8, br_FR.utf8, crh_UA, crh_UA.utf8, cy_GB.utf8, dz_BT, dz_BT.utf8, el_CY.utf8, el_GR.utf8, eo, eo_US.utf8, eo.utf8, et_EE.utf8, eu_ES.utf8, eu_FR.utf8, fa_IR, fa_IR.utf8, fur_IT, fur_IT.utf8, ga_IE.utf8, gd_GB.utf8, gl_ES.utf8, gu_IN, gu_IN.utf8, he_IL.utf8, hi_IN, hi_IN.utf8, hr_HR.utf8, ia_FR, ia_FR.utf8, id_ID.utf8, is_IS.utf8, it_CH.utf8, it_IT.utf8, ja_JP.utf8, ka_GE.utf8, kk_KZ.utf8, km_KH, km_KH.utf8, kn_IN, kn_IN.utf8, ko_KR.utf8, ku_TR.utf8, lt_LT.utf8, lv_LV.utf8, mai_IN, mai_IN.utf8, mk_MK.utf8, ml_IN, ml_IN.utf8, mr_IN, mr_IN.utf8, ms_MY.utf8, my_MM, my_MM.utf8, nb_NO.utf8, nds_DE, nds_DE.utf8, nds_NL, nds_NL.utf8, ne_NP, ne_NP.utf8, nn_NO.utf8, oc_FR.utf8, or_IN, or_IN.utf8, pa_IN, pa_IN.utf8, pa_PK, pa_PK.utf8, pt_BR.utf8, pt_PT.utf8, ro_RO.utf8, ru_RU.utf8, ru_UA.utf8, si_LK, si_LK.utf8, sk_SK.utf8, sl_SI.utf8, sq_AL.utf8, sq_MK, sq_MK.utf8, ta_IN, ta_IN.utf8, ta_LK, ta_LK.utf8, te_IN, te_IN.utf8, tg_TJ.utf8, th_TH.utf8, tr_CY.utf8, tr_TR.utf8, ug_CN, ug_CN.utf8, ug_CN.utf8@latin, ug_CN@latin, xh_ZA.utf8, zh_HK.utf8, zh_SG.utf8, zh_TW.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Estimated minimum size of the filesystem: 797480`,
			expected: 797480,
		},
		{
			name: `hu_HU.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
A fájlrendszer becsült minimális mérete: 797480`,
			expected: 797480,
		},
		{
			name: `da_DK.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Estimeret minimumsstørrelse for filsystemet: 797480`,
			expected: 797480,
		},
		{
			name: `nl_AW, nl_AW.utf8, nl_BE.utf8, nl_NL.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Geschatte minimum grootte van het bestandssysteem: 797480`,
			expected: 797480,
		},
		{
			name: `de_AT.utf8, de_BE.utf8, de_CH.utf8, de_DE.utf8, de_IT.utf8, de_LI.utf8, de_LU.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Geschätzte minimale Größe des Dateisystems: 797480`,
			expected: 797480,
		},
		{
			name: `ca_AD.utf8, ca_ES.utf8, ca_ES.utf8@valencia, ca_FR.utf8, ca_IT.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Mida mínima estimada del sistema de fitxers: 797480`,
			expected: 797480,
		},
		{
			name: `cs_CZ.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Odhadovaná minimální velikost systému souborů: 797480`,
			expected: 797480,
		},
		{
			name: `bs_BA.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Predviđena minimalna veličina datotečnih sistema: 797480`,
			expected: 797480,
		},
		{
			name: `pl_PL.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Przybliżony minimalny rozmiar systemu plików: 797480`,
			expected: 797480,
		},
		{
			name: `fr_BE.utf8, fr_CA.utf8, fr_CH.utf8, fr_FR.utf8, fr_LU.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Taille minimale estimée du système de fichiers : 797480`,
			expected: 797480,
		},
		{
			name: `es_AR.utf8, es_BO.utf8, es_CL.utf8, es_CO.utf8, es_CR.utf8, es_CU, es_CU.utf8, es_DO.utf8, es_EC.utf8, es_ES.utf8, es_GT.utf8, es_HN.utf8, es_MX.utf8, es_NI.utf8, es_PA.utf8, es_PE.utf8, es_PR.utf8, es_PY.utf8, es_SV.utf8, es_US.utf8, es_UY.utf8, es_VE.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Tamaño mínimo estimado del sistema de ficheros: 797480`,
			expected: 797480,
		},
		{
			name: `fi_FI.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Tiedostojärjestelmän arvioitu vähimmäiskoko: 797480`,
			expected: 797480,
		},
		{
			name: `vi_VN, vi_VN.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Ước tính tích cỡ tối thiểu của hệ thống tập tin: 797480`,
			expected: 797480,
		},
		{
			name: `sv_FI.utf8, sv_SE.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Uppskattad minsta storlek på filsystemet: 797480`,
			expected: 797480,
		},
		{
			name: `uk_UA.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Оцінка мінімального розміру файлової системи: 797480`,
			expected: 797480,
		},
		{
			name: `sr_ME, sr_ME.utf8, sr_RS, sr_RS.utf8, sr_RS.utf8@latin, sr_RS@latin`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
Процењена најмања величина система датотека: 797480`,
			expected: 797480,
		},
		{
			name: `zh_CN.utf8`,
			out: `resize2fs 1.44.1 (24-Mar-2018)
预计文件系统的最小尺寸：797480`,
			expected: 797480,
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			minSize, err := parseResize2fsOutputForMinSize(rt.out)
			if err != nil {
				t.Error(err)
			}
			if minSize != rt.expected {
				t.Errorf("expected: %d\n actual: %d", rt.expected, minSize)
			}
		})
	}
}
