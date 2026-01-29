package banner

// GetDisplayAttributes 获取需要展示的属性映射
func (s *Service) GetDisplayAttributes() []map[string]string {
	return []map[string]string{
		{"key": "color", "value": "颜色", "param": "color_id"},
		{"key": "season", "value": "季节", "param": "season_id"},
		{"key": "style", "value": "风格", "param": "style_id"},
		{"key": "occasion", "value": "场合", "param": "occasion_id"},
		{"key": "material", "value": "材质", "param": "material_id"},
		{"key": "collar", "value": "领型", "param": "collar_id"},
		{"key": "sleeve_type", "value": "袖型", "param": "sleeve_type_id"},
		{"key": "sleeve_length", "value": "袖长", "param": "sleeve_length_id"},
		{"key": "gender", "value": "性别", "param": "gender_id"},
		{"key": "clothes_length", "value": "衣长", "param": "clothes_length_id"},
		{"key": "size", "value": "尺码", "param": "size"},
		{"key": "skirt_length", "value": "裙长", "param": "skirt_length_id"},
		{"key": "pattern", "value": "图案", "param": "pattern_id"},
		{"key": "thickness", "value": "厚度", "param": "thickness_id"},
		{"key": "closure", "value": "闭合方式", "param": "closure_id"},
		{"key": "pants_type", "value": "裤型", "param": "pants_type_id"},
		{"key": "pants_length", "value": "裤长", "param": "pants_length_id"},
		{"key": "waistline", "value": "腰围", "param": "waistline_id"},
		{"key": "waist_type", "value": "腰型", "param": "waist_type_id"},
	}
}
