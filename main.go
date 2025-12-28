package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	GraphQLEndpoint = "https://bff-page.kakao.com/graphql"

	// Thay Cookie mới của bạn vào đây
	CookieString = "_kpdid=fe3277a6e854465f872802cdf1e5891c; _kpiid=278f88bc29a85822679573a2aa34ec0f; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kawask=c96dfb31-8e85-4658-867d-e7be7a590f6b; _kau=5d997c1387e3e5ffb5541b773f9197ca5186a653ebf3187d35c5673a32ec93630a0f1153c04e7b3a363c24bd7c0a2ea2262998f6b357f3e9dc9debf5f0555e7623bd02f0b2de4032792563f70285bb611fbba5bbfdba74987764d08064d99375dd7b8c8d92f8e77b17bb92a3799c90f4dc3e97325b0eb3a8d6160ecd31393737393739363530323630393235363234333630343930373235363037330b98995274de95434bab756f125f26bb; _kawlt=egUtYxv4kPI1xNW5C7wEK1eCXjRaKbTrexLvJDM5rdZF8ndiz8cOZdFj9MMc-LBhwFqiTSpdv6rp-9qEXoakiUy9yavVl0ST2R5DvWv6xcAd6BwFkZrPrOdGlyjalnXO; _kawltea=1767004755; _karmt=syalJnnPGHmH1nuNDWZPweABXgAs0-JCl1AQNGu5Gi6IwbhqWrNvJVCjiO1GnvbM; _karmtea=1769531955; _kaslt=6rhhxDwjlVjYBRXQIaaddc+m/C8kp4SjVUs+pSu9KgBL4kdHiQBfyQANuF00vjEo2RnfTO9fagIy9NNYjjgM6A==; _kpawbat_e=pj5hmJUc0MvGobwgdZnUdsxOAJt7kRaBZ%2FSeMQ72ibrFS5AQi7vP%2BsWWryBrS%2FCwzn9M2tLMhhUYNLsYJIGCPEck919ZMeGe3qIfoat2kZ8%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..3RzN1bE0e33WSmjY.NCDd89oPWgMdMFGZqWSIbWirC77EnhtBulLUvCUxx5WYG4rgXfJFYosrUNnmruh22ba-cyIsuql90Tn-s_sA3dA8micrwtMstBqc2im6X2z-RGSWTGlSaBo8WzyW9HkUTrweTNM-wPOXzHYDpa2WNR1Vy-bodF8ABoITPjfNkVYpbOvtQUy0PuD7jPcIXhgjBAA1nlDKWWxaVxuhvrnIE1bgwx6eSgBmtKyonttZznUNGoSh0M1yf3CyM6wDv2V2OgdIPqsTp-SGQeo2c7vxRxV-SBsYK6Gorkigisdqb3RcRL_kEEBo52iUrFT-OtIOZTbLlhFMI5Po6RwPxxFa8Sgsvaw9n5w3HNLQ81Lbr5SUKZJlCUV1ZDafY45i70WIfo3zTY6qWOO2ntZqUGhbAPqgbkeRlKcrP0G2qSm5nHakyVYrnJiKrLLvpKNdIxwvFItoKg1Q4SyJwkOxBc5DUFTKdzbpqrewzJC-V6xm5RQ8Zte32B9kiG3sbbyKk06MA7ADj7N-Yam5oY9mypCrkmFCa7MEDrl6-_gWPYDIv4P4cKCGHfAvz3RbT6sKgRMlZX6QYPYFrD0BUuoMXpdOjBvyEnpSY_HqWbniYX-zTwkJya6QiMOIzXN0lb5xqLu2x1YKSyJik-Dq9Zo9a_x1_Y3q8bcJPPGcEaHySwkusb9L4SweL2Rn5pd1iRJZSImUCOveFxKgRiD4dWATY1qK6_IgG0JF4WaOKHgNvD9E-24ZI3dg1zlBaU0eCFlDdXH7D3dHkhgWDOq1j5J1mn0Ej2gs9ux67vVRmZxpIM-bijVi7wXcv3kAx0Lkkl7B2s1WKBaOv79FllbZzOfebprvO0QkVDRNsRpbpwO5DAxVWr79L48TiPLKIP3hPE1gCZS3zykiPPgHoR02-HEVXdb8LG5Cl4c94bRMOOSRw5ptYBVHG_KbuVpTqT9p4NI6MfXfiYetqmtSwA5JPPwQMVOJ3RwuYPzu_t8KQkPTOR55A5nwQG-NbBvQJGl8yXz1nwBYH3uDdiI_lvXivRwpB8qUztXt2M9JpCVcB3HYYR9stLHhS-mq5SCj_gJni_dyhWh6r9IqNNfrjryA4wyVd9kDeQJS-M_LybMFGvxFpt3iS2WrjHg46tzcbplOJor6SnJ4qLisOuUpxAPuafHZ80mnpJLbZ4CNXZt0aJjUHxzxJiFUOHmdBsdIfaSyk7f-CPvW0lUUOCk0OxPGfxjg7aKI9gm8EJ3qV5yJ2nt0PFTKWmsHgWuLLlypXuzeBzPmdHMVWUv8CJAYO9datVgQz2JGV7x9MqWzqwCL550FDu-I6eSOMganB_HQnEkObIljq-zmUGxqLzlMQ7mMeA91y8JeHyOLJRRpg-0N465UQZjdHi9HEyRoAxJDgJoBDFz0vHjXLVdYnwikK0Ppk10lofJ3Q2tg9jc7uNi-WMx7MGYh4gamBQYzo05hFcSsLel_NSTGa3fBWASsa6NEh7gbW75ubjy2dHktfQUfi9n4CBhyq-MFrIvGSdMMAIiuUHrsa6anGyOvuCOdUBiMs2_ZHhSADs-jH_VVwhErTKRonaUIMaNVUyifFqnKLTYb0Qrg9W2zuErPus6MUwEGvDC9dUZRHmOI0FHfUgag_w8TNAbczowfIAW31ZEx6OAkQYVSF_YJPeFWLNnpR97T4FTLa1ZOZ_GzRt1_totmDMT62Nrt4hPaRUBpLF7i06ooIJ4r7Ly2BJzuTohjPzrk_03eYVbcGE2yXKSD-ZMK0XGBNSAKyFjiqDAqX10OzmtzK65XJaclIYhzj4942Ua9YBwNBePImHgxBm3oQCzFpDLbozVu48SuTc3_HJApmsyHQcR1eKoea93n-EWTDyjgaV5_JWCCCGD1fhLHnmvDoBy160E8oR4r-skg90VFKas3EsqjAfm9hEO39a_PxH2YI7-55Qe-5VrTQ8Vx_2I2l2njjCkiHrxaB5PGvrhde7P3reMgJoH3ldzg987yRuZe4uf_ul9yK45in90f3RHbhfp-adFK3URGPJggXtDD0F2OTWiRxMUsMrNnvNFMcxWNHqWrI_ftjgySQ9wrekaQCILmOKs60T3dIJLrCce40kauXe08m9Qj1UhxRb-k3FfEVCbfAUnJ1E8Aq32lqggl9r-AieeN5h0BMy0uXLRzpsm6EA40chTHFPQga0S2VlEV2_uwFperZ2xo38YU_7rw7ggPyiUrt-GsqDZJEVqHV1KmKLoxEK44BJ6DKKhH6PxzWgp_T4BHQseQ9BKU7ZKcEZjRjxi5SDhz2crD8_EScUZSeytLTu5FlsdYHGw.H_vlm2BWbJlNH75ZpeykkA; _T_ANO=YIy2PEvpPWyqpmGjIT34yxSQCtNyQHqBkuX0kTnugZ0ttKGTIIY7S5q8KMTYEA2jiXQiUeCYbofZ9cP+c0T/ZjDUY3Z3naC36mopRB/gjQKV6xh+OI/VXjPSYK3am7pvZV+gnF0Mg8eX325yhQz1mVzJTwpKXZa9E910Ut9vvLRn9ZfBRtKn8E58rd8siaQ955llklX7U4hJ3uZWMRoeXPE26qE18RPwIXS8uDt5lz0O2NAynAwvw/pi3clKpiB0WgALVARCTlhHcz1ANhha7QJWQbOP0JQLf4jlHaxp8sHAJaLHly+NbRlOM/UPHwhTyZauJOQuLqTzpHMg3/eFIg=="
)

// ============================================================================
// QUERY (dùng payload bạn cung cấp — giữ nguyên các type Long!)
// ============================================================================

const ViewerQuery = `
query viewerInfo($seriesId: Long!, $productId: Long!) {
  viewerInfo(seriesId: $seriesId, productId: $productId) {
    item {
      ...SingleFragment
    }
    seriesItem {
      ...SeriesFragment
    }
    prevItem {
      ...NearItemFragment
    }
    nextItem {
      ...NearItemFragment
    }
    viewerData {
      ...TextViewerData
      ...TalkViewerData
      ...ImageViewerData
      ...VodViewerData
    }
    displayAd {
      ...DisplayAd
    }
  }
}

fragment SingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  series {
    ...SeriesFragment
  }
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  discountRate
  discountRateText
  isShortsDrama
}

fragment SeriesFragment on Series {
  id
  seriesId
  title
  thumbnail
  landThumbnail
  categoryUid
  lang
  category
  categoryType
  subcategoryUid
  subcategory
  badge
  isAllFree
  isWaitfree
  ageGrade
  state
  onIssue
  authors
  description
  pubPeriod
  freeSlideCount
  lastSlideAddedDate
  waitfreeBlockCount
  waitfreePeriodByMinute
  bm
  saleState
  startSaleDt
  saleMethod
  discountRate
  discountRateText
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  translateProperty {
    ...TranslatePropertyFragment
  }
}

fragment ServicePropertyFragment on ServiceProperty {
  viewCount
  readCount
  ratingCount
  ratingSum
  commentCount
  pageContinue {
    ...ContinueInfoFragment
  }
  todayGift {
    ...TodayGift
  }
  preview {
    ...PreviewFragment
  }
  waitfreeTicket {
    ...WaitfreeTicketFragment
  }
  isAlarmOn
  isLikeOn
  ticketCount
  purchasedDate
  lastViewInfo {
    ...LastViewInfoFragment
  }
  purchaseInfo {
    ...PurchaseInfoFragment
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment ContinueInfoFragment on ContinueInfo {
  title
  isFree
  productId
  lastReadProductId
  scheme
  continueProductType
  hasNewSingle
  hasUnreadSingle
}

fragment TodayGift on TodayGift {
  id
  uid
  ticketType
  ticketKind
  ticketCount
  ticketExpireAt
  ticketExpiredText
  isReceived
  seriesId
}

fragment PreviewFragment on Preview {
  item {
    ...PreviewSingleFragment
  }
  nextItem {
    ...PreviewSingleFragment
  }
  usingScroll
}

fragment PreviewSingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
}

fragment OperatorPropertyFragment on OperatorProperty {
  thumbnail
  copy
  helixImpId
  isTextViewer
  selfCensorship
  isBook
  cashInfo {
    discountRate
    setDiscountRate
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment AssetPropertyFragment on AssetProperty {
  bannerImage
  cardImage
  cardTextImage
  cleanImage
  ipxVideo
  bannerSet {
    ...BannerSetFragment
  }
  cardSet {
    ...CardSetFragment
  }
  cardCover {
    ...CardCoverFragment
  }
}

fragment BannerSetFragment on BannerSet {
  backgroundImage
  backgroundColor
  mainImage
  titleImage
}

fragment CardSetFragment on CardSet {
  backgroundColor
  backgroundImage
}

fragment CardCoverFragment on CardCover {
  coverImg
  coverRestricted
}

fragment WaitfreeTicketFragment on WaitfreeTicket {
  chargedPeriod
  chargedCount
  chargedAt
}

fragment LastViewInfoFragment on LastViewInfo {
  isDone
  lastViewDate
  rate
  spineIndex
}

fragment PurchaseInfoFragment on PurchaseInfo {
  purchaseType
  rentExpireDate
  expired
}

fragment TranslatePropertyFragment on TranslateProperty {
  category {
    ...LocaleMapFragment
  }
  sub_category {
    ...LocaleMapFragment
  }
}

fragment LocaleMapFragment on LocaleMap {
  ko
  en
  th
}

fragment NearItemFragment on NearItem {
  productId
  slideType
  ageGrade
  isFree
  title
  thumbnail
}

fragment TextViewerData on TextViewerData {
  type
  atsServerUrl
  metaSecureUrl
  contentsList {
    chapterId
    contentId
    secureUrl
  }
}

fragment TalkViewerData on TalkViewerData {
  type
  talkDownloadData {
    dec
    host
    path
    talkViewerType
  }
}

fragment ImageViewerData on ImageViewerData {
  type
  imageDownloadData {
    ...ImageDownloadData
  }
}

fragment ImageDownloadData on ImageDownloadData {
  files {
    ...ImageDownloadFile
  }
  totalCount
  totalSize
  viewDirection
  gapBetweenImages
  readType
}

fragment ImageDownloadFile on ImageDownloadFile {
  no
  size
  secureUrl
  width
  height
}

fragment VodViewerData on VodViewerData {
  type
  vodDownloadData {
    contentPackId
    drmType
    endpointUrl
    width
    height
    duration
  }
  drmInfo {
    type
    serverType
    error
    fairplayCertificateUrl
    widevineLicenseUrl
    fairplayLicenseUrl
    token
    provider
    assertion
  }
}

fragment DisplayAd on DisplayAd {
  sectionUid
  bannerUid
  treviUid
  momentUid
}
`

// ============================================================================
// MODELS (mình chỉ khai báo những field cần dùng)
// ============================================================================

type GraphQLPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data struct {
		ViewerInfo ViewerInfo `json:"viewerInfo"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ViewerInfo struct {
	Item       Item       `json:"item"`
	SeriesItem *Series    `json:"seriesItem,omitempty"`
	PrevItem   *NearItem  `json:"prevItem,omitempty"`
	NextItem   *NearItem  `json:"nextItem,omitempty"`
	ViewerData ViewerData `json:"viewerData"`
}

type Item struct {
	Title string `json:"title"`
	// thêm field nếu cần về sau
	ProductId int `json:"productId,omitempty"`
	SeriesId  int `json:"seriesId,omitempty"`
}

type Series struct {
	SeriesId int    `json:"seriesId"`
	Title    string `json:"title"`
}

type NearItem struct {
	ProductId int    `json:"productId"`
	Title     string `json:"title"`
	// slideType, etc. omitted — thêm nếu cần
}

type ViewerData struct {
	AtsServerUrl string        `json:"atsServerUrl"`
	MetaSecure   string        `json:"metaSecureUrl,omitempty"`
	ContentsList []ContentPart `json:"contentsList"`
}

type ContentPart struct {
	ChapterId int    `json:"chapterId,omitempty"`
	ContentId int    `json:"contentId,omitempty"`
	SecureUrl string `json:"secureUrl"`
}

type ContentFileResponse struct {
	ContentInfo struct {
		ParagraphList []ParagraphNode `json:"paragraphList"`
	} `json:"contentInfo"`
}

type ParagraphNode struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Children []ParagraphNode `json:"childParagraphList,omitempty"`
}

// ============================================================================
// MAIN: loop theo nextItem.productId — giữ flow tải/chạy y nguyên code gốc
// ============================================================================

func main() {
	// context tổng cho toàn bộ chạy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
		"content-type":    "application/json",
		"origin":          "https://page.kakao.com",
		"referer":         "https://page.kakao.com/",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"cookie":          CookieString,
	}

	// Bắt đầu từ đây: set seriesId và productId chương ban đầu bạn muốn crawl
	currentSeries := 61097897
	currentProduct := 66622589

	for {
		fmt.Printf("\n>>> Crawl productId = %d (seriesId=%d)\n", currentProduct, currentSeries)

		// 1) build payload giữ nguyên như bạn yêu cầu
		payload := GraphQLPayload{
			OperationName: "viewerInfo",
			Query:         ViewerQuery,
			Variables: map[string]interface{}{
				"seriesId":  currentSeries,
				"productId": currentProduct,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		respBody, err := makeRequest(ctx, client, "POST", GraphQLEndpoint, headers, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf(">>> ERROR: GraphQL request failed: %v\n", err)
			break
		}

		var gqlResp GQLResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			fmt.Printf(">>> ERROR: Parse Error: %v\nBody: %s\n", err, string(respBody))
			break
		}
		if len(gqlResp.Errors) > 0 {
			fmt.Printf(">>> ERROR: GraphQL Error: %s\n", gqlResp.Errors[0].Message)
			break
		}

		viewerData := gqlResp.Data.ViewerInfo.ViewerData
		itemInfo := gqlResp.Data.ViewerInfo.Item
		nextItem := gqlResp.Data.ViewerInfo.NextItem

		if viewerData.AtsServerUrl == "" {
			fmt.Println(">>> ERROR: Không tìm thấy atsServerUrl. Kiểm tra lại Cookie.")
			break
		}

		// Tạo file
		safeTitle := sanitizeFilename(itemInfo.Title)
		fileName := safeTitle + ".txt"

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf(">>> ERROR: Không thể tạo file %s: %v\n", fileName, err)
			break
		}
		fmt.Printf(">>> Đã tạo file: %s\n", fileName)
		fmt.Printf(">>> Tổng số phần: %d. Bắt đầu tải...\n", len(viewerData.ContentsList))

		// Ghi tiêu đề
		if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", itemInfo.Title)); err != nil {
			fmt.Printf(">>> WARNING: lỗi ghi tiêu đề: %v\n", err)
		}

		// Loop tải các phần
		for i, part := range viewerData.ContentsList {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}
			fullUrl := viewerData.AtsServerUrl + part.SecureUrl
			fmt.Printf(" -> [%d/%d] Đang tải... (Delay 1s)\n", i+1, len(viewerData.ContentsList))

			partBody, err := makeRequest(ctx, client, "GET", fullUrl, headers, nil)
			if err != nil {
				fmt.Printf("   [Lỗi tải phần %d]: %v\n", i+1, err)
				continue
			}

			var contentResp ContentFileResponse
			if err := json.Unmarshal(partBody, &contentResp); err != nil {
				// nếu không parse được JSON -> có thể là image / binary -> bỏ qua
				continue
			}

			for _, node := range contentResp.ContentInfo.ParagraphList {
				text := extractTextRecursive(node)
				text = strings.ReplaceAll(text, "&nbsp;", " ")
				if _, err := file.WriteString(text); err != nil {
					fmt.Printf("   Lỗi ghi file: %v\n", err)
				}
			}
		}

		file.Close()
		fmt.Println(">>> XONG! Kiểm tra file:", fileName)

		// Nếu không có nextItem -> dừng
		if nextItem == nil || nextItem.ProductId == 0 {
			fmt.Println(">>> HẾT CHƯƠNG hoặc không tìm thấy nextItem.")
			break
		}

		// Cập nhật productId để crawl chương tiếp theo
		fmt.Printf(">>> Next: %d (%s)\n", nextItem.ProductId, nextItem.Title)
		currentProduct = nextItem.ProductId

		// delay an toàn trước lần tiếp theo
		time.Sleep(2 * time.Second)
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func makeRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %d | Body: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}

func extractTextRecursive(node ParagraphNode) string {
	var sb strings.Builder
	if node.Text != "" {
		sb.WriteString(node.Text)
	}
	if node.Type == "BR" {
		sb.WriteString("\n")
	}
	for _, child := range node.Children {
		sb.WriteString(extractTextRecursive(child))
	}
	switch node.Type {
	case "P", "DIV":
		sb.WriteString("\n")
	case "H3":
		sb.WriteString("\n\n### ")
	}
	return sb.String()
}

func sanitizeFilename(name string) string {
	reg, _ := regexp.Compile(`[\\/:*?"<>|]`)
	safe := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(safe)
}
